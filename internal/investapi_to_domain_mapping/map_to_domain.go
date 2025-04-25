package investapi_to_domain_mapping

import (
	"dp/internal"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
)

const (
	divider    = 10000000
	multiplier = 100
)

func MapQuotation(quotation *investapi.Quotation) internal.Percent {
	units := quotation.GetUnits()
	nanos := int64(quotation.GetNano() / divider)
	return internal.Percent(units*multiplier + nanos)
}

func MapMoneyValue(moneyValue *investapi.MoneyValue) internal.Money {
	units := moneyValue.GetUnits()
	nanos := int64(moneyValue.GetNano() / divider)
	return internal.Money(units*multiplier + nanos)
}

func MapQuantity(quotation *investapi.Quotation) internal.Quantity {
	units := quotation.GetUnits()
	return internal.Quantity(units)
}
