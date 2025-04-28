package internal

type PortfolioItem struct {
	Instrument PortfolioInstrument
	Analytics  PortfolioItemAnalytics
}

type PortfolioInstrument struct {
	Instrument Instrument
	Price      Money
	Quantity   Quantity
}

type PortfolioItemAnalytics struct {
	AllTimeMoney   Money
	AllTimePercent Percent
	DailyMoney     Money
	DailyPercent   Percent
	AllMoney       Money
}
