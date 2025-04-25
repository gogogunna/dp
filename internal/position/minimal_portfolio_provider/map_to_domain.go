package minimal_portfolio_provider

import (
	"dp/internal"
	"dp/internal/investapi_to_domain_mapping"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
)

func mapPositions(positions []*investapi.PortfolioPosition) []internal.MinimalPortfolioPosition {
	mappedPositions := make([]internal.MinimalPortfolioPosition, 0, len(positions))
	for _, position := range positions {
		mappedPositions = append(mappedPositions, mapPosition(position))
	}

	return mappedPositions
}

func mapPosition(position *investapi.PortfolioPosition) internal.MinimalPortfolioPosition {
	mappedPosition := internal.Position{
		Figi:     internal.Figi(position.GetFigi()),
		Price:    investapi_to_domain_mapping.MapMoneyValue(position.GetCurrentPrice()),
		Quantity: investapi_to_domain_mapping.MapQuantity(position.GetQuantity()),
	}

	allMoney := mappedPosition.Price * mappedPosition.Quantity
	allTimeMoney := internal.Money(investapi_to_domain_mapping.MapQuotation(position.GetExpectedYield()))
	allTimePercent := int(float64(allTimeMoney) / float64(allMoney) * 10000)
	dailyMoney := investapi_to_domain_mapping.MapMoneyValue(position.GetDailyYield())
	dailyPercent := int(float64(dailyMoney) / float64(allMoney) * 10000)
	return internal.MinimalPortfolioPosition{
		Position:       mappedPosition,
		AllTimeMoney:   allTimeMoney,
		AllTimePercent: allTimePercent,
		DailyMoney:     dailyMoney,
		DailyPercent:   dailyPercent,
		AllMoney:       allMoney,
	}
}
