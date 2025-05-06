package v1

import "time"

type AuthResponse struct {
	OK string `json:"ok"`
}

type MainPageResponse struct {
	Name           string `json:"name"`
	DailyPercent   int    `json:"daily_percent"`
	DailyMoney     int    `json:"daily_money"`
	AlltimePercent int    `json:"alltime_percent"`
	AlltimeMoney   int    `json:"alltime_money"`
	AllMoney       int    `json:"all_money"`
}

type TimeInterval struct {
	From time.Time `json:"from,required"`
	To   time.Time `json:"to,required"`
}

type OperationsRequest struct {
	Figies   []string     `json:"figies,required"`
	Interval TimeInterval `json:"interval,required"`
}

type OperationsResponse struct {
	Operations map[string][]DealOperation `json:"operations_by_figi"`
}

type PortfolioRequest struct {
	Currency int `json:"currency,required"`
}

type PortfolioResponse struct {
	Items          []PortfolioItem `json:"portfolio_items"`
	WarningMessage *string         `json:"warning_message,omitempty"`
}
