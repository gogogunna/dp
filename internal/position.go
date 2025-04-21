package internal

type Figi string

type Position struct {
	Figi     Figi
	Price    Money
	Quantity Quantity
}

type EnrichedPosition struct {
	Position Position
	Name     string
	LogoPath string
}
