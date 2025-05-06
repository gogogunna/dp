package internal

import "dp/pkg/nullable"

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
	AllTimeMoney       Money
	AllTimePercent     Percent
	DailyMoney         Money
	DailyPercent       Percent
	AllMoney           Money
	PercentOfPortfolio Percent
}

type PortfolioAnalytics struct {
	DailyPercent   Percent
	DailyMoney     Money
	AlltimeMoney   Money
	AlltimePercent Percent
	AllMoney       Money
}

type Portfolio struct {
	PortfolioAnalytics PortfolioAnalytics
	Items              []PortfolioItem
	WarningMessage     nullable.Nullable[string]
}
