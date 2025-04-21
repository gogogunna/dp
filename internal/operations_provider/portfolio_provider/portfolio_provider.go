package portfolio_provider

import (
	"context"
	"dp/internal"
	"fmt"
)

type MinimalPortfolioProvider interface {
	Portfolio(
		_ context.Context,
		accountClient internal.AccountIdWithAttachedClientttt,
		currency int,
	) ([]internal.Position, error)
}

type PortfolioProvider struct {
	provider MinimalPortfolioProvider
}

func NewPortfolioProvider(provider MinimalPortfolioProvider) *PortfolioProvider {
	return &PortfolioProvider{
		provider: provider,
	}
}

func (p *PortfolioProvider) Portfolio(
	ctx context.Context,
	accountClient internal.AccountIdWithAttachedClientttt,
	currency int,
) ([]internal.EnrichedPosition, error) {
	portfolio, err := p.provider.Portfolio(ctx, accountClient, currency)
	if err != nil {
		return nil, fmt.Errorf("failed to get minimal portfolio: %w", err)
	}

	figies := make([]internal.EnrichedPosition, len(portfolio))
	return nil, nil
}
