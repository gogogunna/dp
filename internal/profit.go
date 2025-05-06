package internal

type PeriodProfit struct {
	FigiesProfit []FigiProfit
	AllProfit    Money
	Interval     TimeInterval
}

type FigiProfit struct {
	Figi   Figi
	Profit Money
}

type Profit struct {
	Money   Money
	Percent Percent
}
