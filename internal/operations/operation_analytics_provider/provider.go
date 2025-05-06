package operation_analytics_provider

import (
	"context"
	"dp/internal"
	"fmt"
)

const (
	goroutinesMax = 3
)

type OperationsProvider interface {
	DealOperations(
		ctx context.Context,
		accountClient internal.AccountIdWithAttachedClientttt,
		figies []internal.Figi,
		interval internal.TimeInterval,
	) (map[internal.Figi][]internal.DealOperation, error)
}

type OperationAnalyticsCalculator interface {
	CalculateOperationAnalytics(map[internal.OperationType][]internal.DealOperation) internal.OperationAnalyticsItem
}

type OperationAnalyticsProvider struct {
	provider    OperationsProvider
	calculators []OperationAnalyticsCalculator
}

func NewOperationAnalyticsProvider(provider OperationsProvider, calculators ...OperationAnalyticsCalculator) *OperationAnalyticsProvider {
	return &OperationAnalyticsProvider{
		provider:    provider,
		calculators: calculators,
	}
}

func (s *OperationAnalyticsProvider) OperationAnalytics(
	ctx context.Context,
	accountClient internal.AccountIdWithAttachedClientttt,
	interval internal.TimeInterval,
) ([]internal.OperationAnalyticsItem, error) {
	operations, err := s.provider.DealOperations(ctx, accountClient, nil, interval)
	if err != nil {
		return nil, fmt.Errorf("failed to get operations: %w", err)
	}

	operationsByType := make(map[internal.OperationType][]internal.DealOperation, internal.OperationTypesAmount)
	for _, figiOperations := range operations {
		for _, operation := range figiOperations {
			operationsByType[operation.Operation.Type] = append(operationsByType[operation.Operation.Type], operation)
		}
	}

	operationAnalyticsItems := make([]internal.OperationAnalyticsItem, 0, len(s.calculators))
	for _, calculator := range s.calculators {
		answer := calculator.CalculateOperationAnalytics(operationsByType)
		operationAnalyticsItems = append(operationAnalyticsItems, answer)
	}

	return operationAnalyticsItems, nil
}
