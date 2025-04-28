package portfolio_provider

import (
	"dp/internal"
	"dp/internal/client"
	"dp/internal/common_mappers"
)

func MapPortfolioItem(item client.PortfolioItem, instrument client.Instrument) internal.PortfolioItem {
	allMoney := internal.UnitsWithNano(int(item.CurrentPrice) * item.Quantity)
	return internal.PortfolioItem{
		Instrument: internal.PortfolioInstrument{
			Instrument: common_mappers.MapInstrument(instrument),
			Price:      internal.Money(item.CurrentPrice),
			Quantity:   internal.Quantity(item.Quantity),
		},
		Analytics: internal.PortfolioItemAnalytics{
			AllTimeMoney:   internal.Money(item.AllTimeMoney),
			AllTimePercent: internal.CalculatePercent(item.AllTimeMoney, allMoney),
			DailyMoney:     internal.Money(item.DailyMoney),
			DailyPercent:   internal.CalculatePercent(item.DailyMoney, allMoney),
			AllMoney:       internal.Money(allMoney),
		},
	}
}
