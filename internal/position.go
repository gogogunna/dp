package internal

import "dp/pkg/nullable"

type Figi string

type Position struct {
	Figi     Figi
	Price    Money
	Quantity Quantity
}

type MinimalPortfolioPosition struct {
	Position
	AllTimeMoney   Money
	AllTimePercent Percent
	DailyMoney     Money
	DailyPercent   Percent
	AllMoney       Money
}

type PortfolioPosition struct {
	Position MinimalPortfolioPosition
	Enriched PositionEnrichingInfo
}

type PositionEnrichingInfo struct {
	Figi     Figi
	Name     string
	LogoPath string
}

type Operation struct {
	Position
	OperationType        OperationType
	OperationDescription string
}

type EnrichedOperation struct {
	Operation     Operation
	EnrichingInfo nullable.Nullable[PositionEnrichingInfo]
}
