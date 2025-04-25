package main_page

import (
	"context"
	"dp/internal"
	"dp/internal/investapi_to_domain_mapping"
	"fmt"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
)

type MainPageInfoProvider struct{}

func NewMainPageInfoProvider() *MainPageInfoProvider {
	return &MainPageInfoProvider{}
}

func (p *MainPageInfoProvider) MainPageInfo(
	_ context.Context,
	accountClient internal.AccountIdWithAttachedClientttt,
) (internal.MainPageInfo, error) {
	client := accountClient.Client.NewOperationsServiceClient()

	resp, err := client.GetPortfolio(accountClient.AccountId, investapi.PortfolioRequest_RUB)
	if err != nil {
		return internal.MainPageInfo{}, fmt.Errorf("GetPortfolio error: %w", err)
	}

	fmt.Println(accountClient.AccountId)

	portfolioCost := investapi_to_domain_mapping.MapMoneyValue(resp.GetTotalAmountPortfolio())
	percent := investapi_to_domain_mapping.MapQuotation(resp.GetExpectedYield())
	allTimeMoney := int(float64(portfolioCost) * float64(percent) / 10000)
	info := internal.MainPageInfo{
		UserName:       "John Doe",
		DailyPercent:   investapi_to_domain_mapping.MapQuotation(resp.GetDailyYieldRelative()),
		DailyMoney:     investapi_to_domain_mapping.MapMoneyValue(resp.GetDailyYield()),
		AlltimeMoney:   allTimeMoney,
		AlltimePercent: investapi_to_domain_mapping.MapQuotation(resp.GetExpectedYield()),
		AllMoney:       investapi_to_domain_mapping.MapMoneyValue(resp.GetTotalAmountPortfolio()),
	}

	return info, nil
}
