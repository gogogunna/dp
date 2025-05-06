package profit_by_points_calculator

import (
	"context"
	"dp/internal"
	"dp/internal/client"
	"dp/internal/client/instrument_history_price_provider/points_provider"
	"errors"
)

type HistoryPricePointsProvider interface {
	HistoryPrice(
		ctx context.Context,
		accountClient internal.AccountIdWithAttachedClientttt,
		figiesWithPoints []points_provider.FigiWithPoint,
	) ([]client.InstrumentPrice, error)
}

type PointsProfitProvider struct {
	provider HistoryPricePointsProvider
}

func NewPointsProfitProvider(provider HistoryPricePointsProvider) *PointsProfitProvider {
	return &PointsProfitProvider{
		provider: provider,
	}
}

func (s *PointsProfitProvider) Calculate(
	ctx context.Context,
	accountClient internal.AccountIdWithAttachedClientttt,
	figi internal.Figi,
	money internal.Money,
	interval internal.TimeInterval,
) (internal.Profit, error) {
	figiesWithPoint := []points_provider.FigiWithPoint{
		{
			Figi:  figi,
			Point: interval.Start,
		},
		{
			Figi:  figi,
			Point: interval.End,
		},
	}

	prices, err := s.provider.HistoryPrice(ctx, accountClient, figiesWithPoint)
	if err != nil {
		return internal.Profit{}, err
	}

	if len(prices) != 2 {
		return internal.Profit{}, errors.New("expected 2 history prices")
	}
	firstPrice := internal.Money(prices[0].Price)
	secondPrice := internal.Money(prices[1].Price)
	if firstPrice == 0 || secondPrice == 0 {
		return internal.Profit{}, errors.New("prices are empty")
	}

	bought := money / firstPrice
	remainingMoney := money % firstPrice
	profit := remainingMoney + secondPrice*bought
	percent := internal.Percent(float64(secondPrice) / float64(firstPrice))

	return internal.Profit{
		Money:   profit,
		Percent: percent,
	}, nil
}
