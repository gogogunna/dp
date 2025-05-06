package main_page

import (
	"context"
	"dp/internal"
	"fmt"
)

type PortfolioProvider interface {
	PortfolioItems(
		ctx context.Context,
		accountClient internal.AccountIdWithAttachedClientttt,
		currency int,
	) (internal.Portfolio, error)
}

type MainPageInfoProvider struct {
	portfolioProvider PortfolioProvider
}

func NewMainPageInfoProvider(provider PortfolioProvider) *MainPageInfoProvider {
	return &MainPageInfoProvider{
		portfolioProvider: provider,
	}
}

func (p *MainPageInfoProvider) MainPageInfo(
	ctx context.Context,
	accountClient internal.AccountIdWithAttachedClientttt,
) (internal.MainPageInfo, error) {
	resp, err := p.portfolioProvider.PortfolioItems(ctx, accountClient, 0)
	if err != nil {
		return internal.MainPageInfo{}, fmt.Errorf("failed to get porfolio: %w", err)
	}

	info := internal.MainPageInfo{
		UserName:           "John Doe",
		PortfolioAnalytics: resp.PortfolioAnalytics,
	}

	return info, nil
}
