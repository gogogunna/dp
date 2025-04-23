package internal

import "dp/pkg/nullable"

type Figi string

type Position struct {
	Figi     Figi
	Price    Money
	Quantity Quantity
}

type PortfolioPosition struct {
	Position Position
	Enriched PositionEnrichingInfo
}

type PositionEnrichingInfo struct {
	Figi     Figi
	Name     string
	LogoPath string
}

type Operation struct {
	Position
	OperationType OperationType
}

type EnrichedOperation struct {
	Operation     Operation
	EnrichingInfo nullable.Nullable[PositionEnrichingInfo]
}
