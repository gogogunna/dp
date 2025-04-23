package operations_provider

import (
	"dp/internal"
	"dp/internal/investapi_to_domain_mapping"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
)

func mapOperationToDomain(resp *investapi.Operation) internal.Operation {
	return internal.Operation{
		Position: internal.Position{
			Figi:     internal.Figi(resp.GetFigi()),
			Price:    investapi_to_domain_mapping.MapMoneyValue(resp.GetPrice()),
			Quantity: internal.Quantity(resp.GetQuantity()),
		},
		OperationType: internal.OperationType(resp.GetOperationType()),
	}
}

func mapOperationsToDomain(resp *investapi.OperationsResponse) []internal.Operation {
	ops := make([]internal.Operation, 0, len(resp.Operations))
	for _, op := range resp.Operations {
		ops = append(ops, mapOperationToDomain(op))
	}

	return ops
}
