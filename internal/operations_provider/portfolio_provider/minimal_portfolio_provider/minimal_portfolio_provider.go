package minimal_portfolio_provider

import (
	"context"
	"dp/internal"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
)

type MinimalPortfolioProvider struct{}

func NewOperationsProvider() *MinimalPortfolioProvider {
	return &MinimalPortfolioProvider{}
}

func (p *MinimalPortfolioProvider) Portfolio(
	_ context.Context,
	accountClient internal.AccountIdWithAttachedClientttt,
	currency int,
) ([]internal.Position, error) {
	operationsClient := accountClient.Client.NewOperationsServiceClient()

	portfolio, err := operationsClient.GetPortfolio(accountClient.AccountId, investapi.PortfolioRequest_CurrencyRequest(currency))
	if err != nil {
		return nil, err
	}

	return mapPositions(portfolio.GetPositions()), nil
}
