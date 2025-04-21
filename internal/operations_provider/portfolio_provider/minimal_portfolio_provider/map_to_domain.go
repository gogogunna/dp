package minimal_portfolio_provider

import (
	"dp/internal"
	"dp/internal/investapi_to_domain_mapping"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
)

func mapPositions(positions []*investapi.PortfolioPosition) []internal.Position {
	mappedPositions := make([]internal.Position, 0, len(positions))
	for _, position := range positions {
		mappedPositions = append(mappedPositions, mapPosition(position))
	}

	return mappedPositions
}

func mapPosition(position *investapi.PortfolioPosition) internal.Position {
	return internal.Position{
		Figi:     internal.Figi(position.Figi),
		Price:    investapi_to_domain_mapping.MapMoneyValue(position.AveragePositionPrice),
		Quantity: investapi_to_domain_mapping.MapQuantity(position.Quantity),
	}
}
