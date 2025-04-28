package v1

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

type OperationsRequest struct {
	Figies   []string     `json:"figies"`
	Interval TimeInterval `json:"interval"`
}

type OperationsResponse struct {
	Operations map[string][]DealOperation `json:"operations_by_figi"`
}

type PortfolioRequest struct {
	Currency int `json:"currency"`
}

type PortfolioResponse struct {
	Items []PortfolioItem `json:"portfolio_items"`
}
