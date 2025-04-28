package portfolio_provider

import (
	"context"
	"dp/internal"
	"dp/internal/client"
	"fmt"
)

type PortfolioProvider interface {
	Portfolio(
		_ context.Context,
		accountClient internal.AccountIdWithAttachedClientttt,
		currency int,
	) (client.Portfolio, error)
}

type InstrumentsInfoProvider interface {
	InstrumentsInfo(
		_ context.Context,
		accountClient internal.AccountIdWithAttachedClientttt,
		figies []internal.Figi,
	) (map[internal.Figi]client.Instrument, error)
}

type PortfolioItemsProvider struct {
	portfolioProvider       PortfolioProvider
	instrumentsInfoProvider InstrumentsInfoProvider
}

func NewPortfolioItemsProvider(
	portfolioProvider PortfolioProvider,
	instrumentsInfoProvider InstrumentsInfoProvider,
) *PortfolioItemsProvider {
	return &PortfolioItemsProvider{
		portfolioProvider:       portfolioProvider,
		instrumentsInfoProvider: instrumentsInfoProvider,
	}
}

func (p *PortfolioItemsProvider) PortfolioItems(
	ctx context.Context,
	accountClient internal.AccountIdWithAttachedClientttt,
	currency int,
) ([]internal.PortfolioItem, error) {
	portfolio, err := p.portfolioProvider.Portfolio(ctx, accountClient, currency)
	if err != nil {
		return nil, fmt.Errorf("failed to get minimal portfolio: %w", err)
	}

	figies := make([]internal.Figi, 0, len(portfolio.Items))
	for _, position := range portfolio.Items {
		figies = append(figies, internal.Figi(position.Figi))
	}

	instrumentsInfo, err := p.instrumentsInfoProvider.InstrumentsInfo(ctx, accountClient, figies)
	if err != nil {
		return nil, fmt.Errorf("failed to get position enriching info: %w", err)
	}

	portfolioItems := make([]internal.PortfolioItem, 0, len(instrumentsInfo))
	for _, item := range portfolio.Items {

		portfolioItems = append(portfolioItems, MapPortfolioItem(item, instrumentsInfo[internal.Figi(item.Figi)]))
	}

	return portfolioItems, nil
}
