package operations_provider

import (
	"context"
	"dp/internal"
	"dp/internal/client"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
)

type OperationsProvider struct{}

func NewOperationsProvider() *OperationsProvider {
	return &OperationsProvider{}
}

func (p *OperationsProvider) Operations(
	_ context.Context,
	accountClient internal.AccountIdWithAttachedClientttt,
	figies []internal.Figi,
	interval internal.TimeInterval,
) ([]client.Operation, error) {
	operationsClient := accountClient.Client.NewOperationsServiceClient()

	figi := ""
	if len(figies) == 1 {
		figi = string(figies[0])
	}
	req := &investgo.GetOperationsRequest{
		Figi:      figi,
		AccountId: accountClient.AccountId,
		State:     investapi.OperationState_OPERATION_STATE_EXECUTED,
		From:      interval.Start,
		To:        interval.End,
	}
	operationsResponse, err := operationsClient.GetOperations(req)
	if err != nil {
		return nil, err
	}

	if len(figies) != 0 {
		figiesHashSet := make(map[string]struct{}, len(figies))
		for _, fig := range figies {
			figiesHashSet[string(fig)] = struct{}{}
		}

		operations := make([]*investapi.Operation, 0, len(operationsResponse.Operations))
		for _, operation := range operationsResponse.GetOperations() {
			if _, ok := figiesHashSet[operation.GetFigi()]; ok {
				operations = append(operations, operation)
			}
		}

		return mapOperationsToDomain(operations), nil
	}

	return mapOperationsToDomain(operationsResponse.OperationsResponse.GetOperations()), nil
}
