package buy

import (
	"dp/internal"
	"dp/pkg/nullable"
)

type BuyCalculator struct{}

func NewBuyCalculator() *BuyCalculator {
	return &BuyCalculator{}
}

func (f *BuyCalculator) CalculateOperationAnalytics(operationsByType map[internal.OperationType][]internal.DealOperation) internal.OperationAnalyticsItem {
	moneyValue := internal.Money(0)
	for _, operation := range operationsByType[internal.OperationTypeBuy] {
		moneyValue += operation.Operation.Payment.Value()
	}

	return internal.OperationAnalyticsItem{
		Number:     1,
		Text:       "Потрачено на комиссии",
		MoneySpent: nullable.NewValue(moneyValue),
	}
}
