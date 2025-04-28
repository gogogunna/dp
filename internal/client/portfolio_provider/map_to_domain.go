package portfolio_provider

import (
	"dp/internal"
	"dp/internal/client"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
)

func mapPortfolio(resp *investgo.PortfolioResponse) client.Portfolio {
	return client.Portfolio{
		AllMoney:       internal.MapUnitsWithNano(resp.GetTotalAmountPortfolio()),
		AlltimePercent: internal.MapUnitsWithNano(resp.GetExpectedYield()),
		DailyPercent:   internal.MapUnitsWithNano(resp.GetDailyYieldRelative()),
		DailyMoney:     internal.MapUnitsWithNano(resp.GetDailyYield()),
		Items:          mapPositions(resp.GetPositions()),
	}
}

func mapPositions(positions []*investapi.PortfolioPosition) []client.PortfolioItem {
	mappedPositions := make([]client.PortfolioItem, 0, len(positions))
	for _, position := range positions {
		mappedPositions = append(mappedPositions, mapPosition(position))
	}

	return mappedPositions
}

func mapPosition(position *investapi.PortfolioPosition) client.PortfolioItem {
	return client.PortfolioItem{
		Figi:         position.GetFigi(),
		CurrentPrice: internal.MapUnitsWithNano(position.GetCurrentPrice()),
		Quantity:     int(position.GetQuantity().GetUnits()),
		AllTimeMoney: internal.MapUnitsWithNano(position.GetExpectedYield()),
		DailyMoney:   internal.MapUnitsWithNano(position.GetDailyYield()),
	}
}
