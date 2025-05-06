package portfolio_provider

import (
	"context"
	"dp/internal"
	"dp/internal/client"
	"dp/pkg/nullable"
	"fmt"
)

type ClientPortfolioProvider interface {
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

type PortfolioProvider struct {
	portfolioProvider       ClientPortfolioProvider
	instrumentsInfoProvider InstrumentsInfoProvider
}

func NewPortfolioProvider(
	portfolioProvider ClientPortfolioProvider,
	instrumentsInfoProvider InstrumentsInfoProvider,
) *PortfolioProvider {
	return &PortfolioProvider{
		portfolioProvider:       portfolioProvider,
		instrumentsInfoProvider: instrumentsInfoProvider,
	}
}

func (p *PortfolioProvider) PortfolioItems(
	ctx context.Context,
	accountClient internal.AccountIdWithAttachedClientttt,
	currency int,
) (internal.Portfolio, error) {
	portfolio, err := p.portfolioProvider.Portfolio(ctx, accountClient, currency)
	if err != nil {
		return internal.Portfolio{}, fmt.Errorf("failed to get minimal portfolio: %w", err)
	}

	figies := make([]internal.Figi, 0, len(portfolio.Items))
	for _, position := range portfolio.Items {
		figies = append(figies, internal.Figi(position.Figi))
	}

	instrumentsInfo, err := p.instrumentsInfoProvider.InstrumentsInfo(ctx, accountClient, figies)
	if err != nil {
		return internal.Portfolio{}, fmt.Errorf("failed to get position enriching info: %w", err)
	}

	portfolioAnalytics := mapPortfolioAnalytics(portfolio)
	warningMessage := nullable.Nullable[string]{}
	portfolioItems := make([]internal.PortfolioItem, 0, len(instrumentsInfo))
	for _, item := range portfolio.Items {
		portfolioItem, message := MapPortfolioItem(portfolio.AllMoney, item, instrumentsInfo[internal.Figi(item.Figi)])
		if message.IsDefined() {
			warningMessage = message
		}

		portfolioItems = append(portfolioItems, portfolioItem)
	}

	return internal.Portfolio{
		PortfolioAnalytics: portfolioAnalytics,
		Items:              portfolioItems,
		WarningMessage:     warningMessage,
	}, nil
}
