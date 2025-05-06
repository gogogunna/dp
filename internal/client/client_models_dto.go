package client

import (
	"dp/internal"
	"time"
)

type Operation struct {
	Figi          string
	OperationType int
	Time          time.Time
	Payment       internal.UnitsWithNano
	Quantity      int
	Price         internal.UnitsWithNano
}

type Instrument struct {
	Figi     string
	Name     string
	LogoPath string
}

type Portfolio struct {
	AllMoney       internal.UnitsWithNano
	AlltimePercent internal.UnitsWithNano
	DailyPercent   internal.UnitsWithNano
	DailyMoney     internal.UnitsWithNano
	Items          []PortfolioItem
}

type PortfolioItem struct {
	Figi         string
	CurrentPrice internal.UnitsWithNano
	Quantity     int
	AllTimeMoney internal.UnitsWithNano
	DailyMoney   internal.UnitsWithNano
}

type InstrumentPrice struct {
	Price     internal.UnitsWithNano
	RealTime  time.Time
	PointTime time.Time
}
