package v1

import (
	"dp/internal"
	"dp/pkg/nullable"
)

func mapInstrument(instrument internal.Instrument) Instrument {
	return Instrument{
		Figi:     string(instrument.Figi),
		Name:     instrument.Name,
		LogoPath: instrument.LogoPath,
	}
}

func mapDealOperation(dealOperation internal.DealOperation) DealOperation {
	return DealOperation{
		Operation: mapOperation(dealOperation.Operation),
		Deal:      mapOperationDeal(dealOperation.Deal),
	}
}

func mapOperation(operation internal.Operation) Operation {
	payment := 0
	var paymentPtr *int
	if operation.Payment.IsDefined() {
		payment = int(operation.Payment.Value())
		paymentPtr = &payment
	}
	return Operation{
		Type:        int(operation.Type),
		Description: operation.Description,
		Time:        operation.Time,
		Payment:     paymentPtr,
	}
}

func mapOperationDeal(operationDeal nullable.Nullable[internal.Deal]) *OperationDeal {
	if operationDeal.IsNil() {
		return nil
	}
	deal := operationDeal.Value()
	return &OperationDeal{
		Instrument: mapInstrument(deal.Instrument),
		Price:      int(deal.Price),
		Quantity:   int(deal.Quantity),
	}
}

func mapPortfolioItem(item internal.PortfolioItem) PortfolioItem {
	return PortfolioItem{
		Instrument: mapPortfolioInstrument(item.Instrument),
		Analytics:  mapPortfolioAnalytics(item.Analytics),
	}
}

func mapPortfolioInstrument(item internal.PortfolioInstrument) PortfolioInstrument {
	return PortfolioInstrument{
		Instrument: mapInstrument(item.Instrument),
		Price:      int(item.Price),
		Quantity:   int(item.Quantity),
	}
}

func mapPortfolioAnalytics(analytics internal.PortfolioItemAnalytics) PortfolioInstrumentAnalytics {
	return PortfolioInstrumentAnalytics{
		AllTimeMoney:       int(analytics.AllTimeMoney),
		AllTimePercent:     int(analytics.AllTimePercent),
		DailyMoney:         int(analytics.DailyMoney),
		DailyPercent:       int(analytics.DailyPercent),
		AllMoney:           int(analytics.AllMoney),
		PercentOfPortfolio: int(analytics.PercentOfPortfolio),
	}
}
