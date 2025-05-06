package deal_operation_provider

import (
	"dp/internal"
	"dp/internal/client"
	"dp/internal/common_mappers"
	"dp/pkg/nullable"
)

func mapOperation(operation client.Operation) internal.Operation {
	return internal.Operation{
		Type:        internal.OperationType(operation.OperationType),
		Description: internal.OperationTypeDescs[internal.OperationType(operation.OperationType)],
		Time:        operation.Time,
		Payment:     nullable.NilDefaultValuee(internal.Money(operation.Payment)),
	}
}

func mapDeal(operation client.Operation, instrumentInfo client.Instrument) internal.Deal {
	return internal.Deal{
		Instrument: common_mappers.MapInstrument(instrumentInfo),
		Price:      internal.Money(operation.Price),
		Quantity:   internal.Quantity(operation.Quantity),
	}
}

func mapDealOperation(operation client.Operation, instrumentInfo client.Instrument) internal.DealOperation {
	return internal.DealOperation{
		Operation: mapOperation(operation),
		Deal:      nullable.NilDefaultValuee(mapDeal(operation, instrumentInfo)),
	}
}

func mapOperationWithoutDeal(operation client.Operation) internal.DealOperation {
	return internal.DealOperation{
		Operation: mapOperation(operation),
	}
}
