package v1

import (
	"time"
)

type DealOperation struct {
	Operation Operation      `json:"operation"`
	Deal      *OperationDeal `json:"deal"`
}

type Operation struct {
	Type        int       `json:"type"`
	Description string    `json:"description"`
	Time        time.Time `json:"time"`
	Payment     *int      `json:"payment"`
}

type OperationDeal struct {
	Instrument Instrument `json:"instrument"`
	Price      int        `json:"price"`
	Quantity   int        `json:"quantity"`
}

type Instrument struct {
	Figi     string `json:"figi"`
	Name     string `json:"name"`
	LogoPath string `json:"logo_path"`
}

type PortfolioInstrument struct {
	Instrument Instrument `json:"instrument"`
	Price      int        `json:"price"`
	Quantity   int        `json:"quantity"`
}

type PortfolioInstrumentAnalytics struct {
	AllTimeMoney   int `json:"all_time_money"`
	AllTimePercent int `json:"all_time_percent"`
	DailyMoney     int `json:"daily_money"`
	DailyPercent   int `json:"daily_percent"`
	AllMoney       int `json:"all_money"`
}

type PortfolioItem struct {
	Instrument PortfolioInstrument          `json:"portfolio_instrument"`
	Analytics  PortfolioInstrumentAnalytics `json:"portfolio_instrument_analytics"`
}

type TimeInterval struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}
