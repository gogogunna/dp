package enriched_operations_provider

import (
	"context"
	"dp/internal"
	"dp/pkg/nullable"
	"dp/pkg/slices"
	"fmt"
)

type OperationsProvider interface {
	Operations(
		_ context.Context,
		accountClient internal.AccountIdWithAttachedClientttt,
		figies []internal.Figi,
		interval internal.TimeInterval,
	) ([]internal.Operation, error)
}

type PositionEnrichingInfoProvider interface {
	PositionEnrichingInfo(
		_ context.Context,
		accountClient internal.AccountIdWithAttachedClientttt,
		figies []internal.Figi,
	) (map[internal.Figi]internal.PositionEnrichingInfo, error)
}

type EnrichedOperationsProvider struct {
	operationsProvider    OperationsProvider
	positionEnrichingInfo PositionEnrichingInfoProvider
}

func NewEnrichedOperationsProvider(
	operationsProvider OperationsProvider,
	infoProvider PositionEnrichingInfoProvider,
) *EnrichedOperationsProvider {
	return &EnrichedOperationsProvider{
		operationsProvider:    operationsProvider,
		positionEnrichingInfo: infoProvider,
	}
}

func (p *EnrichedOperationsProvider) Operations(
	ctx context.Context,
	accountClient internal.AccountIdWithAttachedClientttt,
	figies []internal.Figi,
	interval internal.TimeInterval,
) (map[internal.Figi][]internal.EnrichedOperation, error) {
	operations, err := p.operationsProvider.Operations(ctx, accountClient, figies, interval)
	if err != nil {
		return nil, fmt.Errorf("failed to get minimal portfolio: %w", err)
	}

	operationFigies := make([]internal.Figi, 0, len(figies))
	for _, operation := range operations {
		if operation.Figi != "" {
			operationFigies = append(operationFigies, operation.Figi)
		}
	}

	operationFigies = slices.Unique(operationFigies)

	enrichingInfo, err := p.positionEnrichingInfo.PositionEnrichingInfo(ctx, accountClient, operationFigies)
	if err != nil {
		return nil, fmt.Errorf("failed to get position enriching info: %w", err)
	}

	enrichedOperations := make(map[internal.Figi][]internal.EnrichedOperation, len(enrichingInfo))
	for _, item := range operations {
		enriched := internal.EnrichedOperation{
			Operation: item,
		}

		if item.Figi != "" {
			enriching, ok := enrichingInfo[item.Figi]
			if ok {
				enriched.EnrichingInfo = nullable.NewValue(enriching)
			}
		}

		figi := item.Figi
		if len(figies) == 0 {
			figi = ""
		}
		enrichedOperations[figi] = append(
			enrichedOperations[figi],
			enriched,
		)
	}

	return enrichedOperations, nil
}
