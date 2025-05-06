package sell

import (
	"dp/internal"
	"dp/pkg/nullable"
)

type SellCalculator struct{}

func NewSellCalculator() *SellCalculator {
	return &SellCalculator{}
}

func (f *SellCalculator) CalculateOperationAnalytics(operationsByType map[internal.OperationType][]internal.DealOperation) internal.OperationAnalyticsItem {
	moneyValue := internal.Money(0)
	for _, operation := range operationsByType[internal.OperationTypeSell] {
		moneyValue += operation.Operation.Payment.Value()
	}

	return internal.OperationAnalyticsItem{
		Number:     1,
		Text:       "Потрачено на комиссии",
		MoneySpent: nullable.NewValue(moneyValue),
	}
}
