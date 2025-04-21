package investapi_to_domain_mapping

import (
	"dp/internal"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
)

const (
	nanoDivider = 1000000000
)

func MapQuotation(quotation *investapi.Quotation) internal.Percent {
	units := quotation.GetUnits()
	nanos := int64(float64(quotation.GetNano()) / nanoDivider)
	return internal.Percent(units + nanos)
}

func MapMoneyValue(moneyValue *investapi.MoneyValue) internal.Money {
	units := moneyValue.GetUnits()
	nanos := int64(float64(moneyValue.GetNano()) / nanoDivider)
	return internal.Money(units + nanos)
}

func MapQuantity(quotation *investapi.Quotation) internal.Percent {
	units := quotation.GetUnits()
	nanos := int64(float64(quotation.GetNano()) / nanoDivider)
	return internal.Percent(units + nanos)
}
