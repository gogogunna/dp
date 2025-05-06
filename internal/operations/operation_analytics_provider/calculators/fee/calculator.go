package fee

import (
	"dp/internal"
	"dp/pkg/nullable"
)

var neededOperationTypes = []internal.OperationType{
	internal.OperationTypeAdviceFee,
	internal.OperationTypeBrokerFee,
	internal.OperationTypeCashFee,
	internal.OperationTypeMarginFee,
	internal.OperationTypeOutFee,
	internal.OperationTypeSuccessFee,
}

type FeeCalculator struct{}

func NewFeeCalculator() *FeeCalculator {
	return &FeeCalculator{}
}

func (f *FeeCalculator) CalculateOperationAnalytics(operationsByType map[internal.OperationType][]internal.DealOperation) internal.OperationAnalyticsItem {
	moneyValue := internal.Money(0)
	for _, opType := range neededOperationTypes {
		if operations, ok := operationsByType[opType]; ok {
			for _, operation := range operations {
				moneyValue += operation.Operation.Payment.Value()
			}
		}
	}

	return internal.OperationAnalyticsItem{
		Number:     1,
		Text:       "Потрачено на комиссии",
		MoneySpent: nullable.NewValue(moneyValue),
	}
}
