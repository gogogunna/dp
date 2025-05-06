package portfolio_provider

import (
	"dp/internal"
	"dp/internal/client"
	"dp/internal/common_mappers"
	"dp/pkg/nullable"
	"fmt"
)

const percentBorder = 80

func mapPortfolioAnalytics(portfolio client.Portfolio) internal.PortfolioAnalytics {
	return internal.PortfolioAnalytics{
		DailyPercent:   internal.Percent(portfolio.DailyPercent),
		DailyMoney:     internal.Money(portfolio.DailyMoney),
		AlltimeMoney:   internal.CalculateDiff(portfolio.DailyPercent, portfolio.DailyMoney),
		AlltimePercent: internal.Percent(portfolio.AlltimePercent),
		AllMoney:       internal.Money(portfolio.AllMoney),
	}
}

func MapPortfolioItem(portfolioMoney internal.UnitsWithNano, item client.PortfolioItem, instrument client.Instrument) (internal.PortfolioItem, nullable.Nullable[string]) {
	warningMessage := nullable.Nullable[string]{}
	allMoney := internal.UnitsWithNano(int(item.CurrentPrice) * item.Quantity)
	portfolioItem := internal.PortfolioItem{
		Instrument: internal.PortfolioInstrument{
			Instrument: common_mappers.MapInstrument(instrument),
			Price:      internal.Money(item.CurrentPrice),
			Quantity:   internal.Quantity(item.Quantity),
		},
		Analytics: internal.PortfolioItemAnalytics{
			AllTimeMoney:       internal.Money(item.AllTimeMoney),
			AllTimePercent:     internal.CalculatePercent(item.AllTimeMoney, allMoney),
			DailyMoney:         internal.Money(item.DailyMoney),
			DailyPercent:       internal.CalculatePercent(item.DailyMoney, allMoney),
			AllMoney:           internal.Money(allMoney),
			PercentOfPortfolio: internal.CalculatePercent(allMoney, portfolioMoney),
		},
	}

	if portfolioItem.Analytics.PercentOfPortfolio > percentBorder*100 {
		warningMessage.SetValue(fmt.Sprintf("В вашем портфеле перекос. %s составляет более %d от портфеля.", portfolioItem.Instrument.Instrument.Name, percentBorder))
	}

	return portfolioItem, warningMessage
}
