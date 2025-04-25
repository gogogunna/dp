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
	) ([]internal.MinimalPortfolioPosition, error)
}

type PositionEnrichingInfoProvider interface {
	PositionEnrichingInfo(
		_ context.Context,
		accountClient internal.AccountIdWithAttachedClientttt,
		figies []internal.Figi,
	) (map[internal.Figi]internal.PositionEnrichingInfo, error)
}

type PortfolioProvider struct {
	portfolioProvider     MinimalPortfolioProvider
	positionEnrichingInfo PositionEnrichingInfoProvider
}

func NewPortfolioProvider(
	portfolioProvider MinimalPortfolioProvider,
	infoProvider PositionEnrichingInfoProvider,
) *PortfolioProvider {
	return &PortfolioProvider{
		portfolioProvider:     portfolioProvider,
		positionEnrichingInfo: infoProvider,
	}
}

func (p *PortfolioProvider) Portfolio(
	ctx context.Context,
	accountClient internal.AccountIdWithAttachedClientttt,
	currency int,
) ([]internal.PortfolioPosition, error) {
	portfolio, err := p.portfolioProvider.Portfolio(ctx, accountClient, currency)
	if err != nil {
		return nil, fmt.Errorf("failed to get minimal portfolio: %w", err)
	}

	figies := make([]internal.Figi, 0, len(portfolio))
	for _, position := range portfolio {
		figies = append(figies, position.Figi)
	}

	enrichingInfo, err := p.positionEnrichingInfo.PositionEnrichingInfo(ctx, accountClient, figies)
	if err != nil {
		return nil, fmt.Errorf("failed to get position enriching info: %w", err)
	}

	enrichedPositions := make([]internal.PortfolioPosition, 0, len(enrichingInfo))
	for _, item := range portfolio {
		enrichedPositions = append(enrichedPositions, internal.PortfolioPosition{
			Position: item,
			Enriched: enrichingInfo[item.Figi],
		})
	}

	return enrichedPositions, nil
}
