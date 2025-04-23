package operations_provider

import (
	"context"
	"dp/internal"
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
) ([]internal.Operation, error) {
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
	portfolio, err := operationsClient.GetOperations(req)
	if err != nil {
		return nil, err
	}

	return mapOperationsToDomain(portfolio.OperationsResponse), nil
}
