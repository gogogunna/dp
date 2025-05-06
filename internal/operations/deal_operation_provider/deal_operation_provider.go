package deal_operation_provider

import (
	"context"
	"dp/internal"
	"dp/internal/client"
	"dp/pkg/slices"
	"fmt"
)

type OperationsProvider interface {
	Operations(
		_ context.Context,
		accountClient internal.AccountIdWithAttachedClientttt,
		figies []internal.Figi,
		interval internal.TimeInterval,
	) ([]client.Operation, error)
}

type InstrumentsInfoProvider interface {
	InstrumentsInfo(
		_ context.Context,
		accountClient internal.AccountIdWithAttachedClientttt,
		figies []internal.Figi,
	) (map[internal.Figi]client.Instrument, error)
}

type DealOperationsProvider struct {
	operationsProvider      OperationsProvider
	instrumentsInfoProvider InstrumentsInfoProvider
}

func NewDealOperationsProvider(
	operationsProvider OperationsProvider,
	instrumentsInfoProvider InstrumentsInfoProvider,
) *DealOperationsProvider {
	return &DealOperationsProvider{
		operationsProvider:      operationsProvider,
		instrumentsInfoProvider: instrumentsInfoProvider,
	}
}

func (p *DealOperationsProvider) DealOperations(
	ctx context.Context,
	accountClient internal.AccountIdWithAttachedClientttt,
	figies []internal.Figi,
	interval internal.TimeInterval,
) (map[internal.Figi][]internal.DealOperation, error) {
	operations, err := p.operationsProvider.Operations(ctx, accountClient, figies, interval)
	if err != nil {
		return nil, fmt.Errorf("failed to get minimal portfolio: %w", err)
	}

	operationFigies := make([]internal.Figi, 0, len(figies))
	for _, operation := range operations {
		if operation.Figi != "" {
			operationFigies = append(operationFigies, internal.Figi(operation.Figi))
		}
	}

	operationFigies = slices.Unique(operationFigies)

	instrumentsInfo, err := p.instrumentsInfoProvider.InstrumentsInfo(ctx, accountClient, operationFigies)
	if err != nil {
		return nil, fmt.Errorf("failed to get position enriching info: %w", err)
	}

	dealOperations := make(map[internal.Figi][]internal.DealOperation, len(instrumentsInfo))
	for _, operation := range operations {
		var dealOperation internal.DealOperation

		if operation.Figi != "" {
			instrument, ok := instrumentsInfo[internal.Figi(operation.Figi)]
			if ok {
				dealOperation = mapDealOperation(operation, instrument)
			}
		} else {
			dealOperation = mapOperationWithoutDeal(operation)
		}

		figi := internal.Figi(operation.Figi)
		if len(figies) == 0 {
			figi = ""
		}
		dealOperations[figi] = append(
			dealOperations[figi],
			dealOperation,
		)
	}

	return dealOperations, nil
}
