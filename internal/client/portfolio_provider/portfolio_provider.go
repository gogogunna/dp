package portfolio_provider

import (
	"context"
	"dp/internal"
	"dp/internal/client"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
)

type PortfolioProvider struct{}

func NewPortfolioProvider() *PortfolioProvider {
	return &PortfolioProvider{}
}

func (p *PortfolioProvider) Portfolio(
	_ context.Context,
	accountClient internal.AccountIdWithAttachedClientttt,
	currency int,
) (client.Portfolio, error) {
	operationsClient := accountClient.Client.NewOperationsServiceClient()

	portfolio, err := operationsClient.GetPortfolio(accountClient.AccountId, investapi.PortfolioRequest_CurrencyRequest(currency))
	if err != nil {
		return client.Portfolio{}, err
	}

	return mapPortfolio(portfolio), nil
}
