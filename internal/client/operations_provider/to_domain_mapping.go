package operations_provider

import (
	"dp/internal"
	"dp/internal/client"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
)

func mapOperationToDomain(resp *investapi.Operation) client.Operation {
	return client.Operation{
		Figi:          resp.GetFigi(),
		OperationType: int(resp.GetOperationType()),
		Time:          resp.GetDate().AsTime(),
		Payment:       internal.MapUnitsWithNano(resp.GetPayment()),
		Quantity:      int(resp.GetQuantity()),
		Price:         internal.MapUnitsWithNano(resp.GetPrice()),
	}
}

func mapOperationsToDomain(operations []*investapi.Operation) []client.Operation {
	ops := make([]client.Operation, 0, len(operations))
	for _, op := range operations {
		ops = append(ops, mapOperationToDomain(op))
	}

	return ops
}
