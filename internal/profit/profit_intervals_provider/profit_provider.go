package profit_intervals_provider

import (
	"context"
	"dp/internal"
	"dp/internal/client"
	"fmt"
	"time"
)

type HistoryPriceBatchProvider interface {
	HistoryPriceBatch(
		ctx context.Context,
		accountClient internal.AccountIdWithAttachedClientttt,
		figies []internal.Figi,
		interval internal.TimeInterval,
		points int,
	) (map[internal.Figi][]client.InstrumentPrice, error)
}

type PortfolioProvider interface {
	Portfolio(
		_ context.Context,
		accountClient internal.AccountIdWithAttachedClientttt,
		currency int,
	) (client.Portfolio, error)
}

type OperationsProvider interface {
	Operations(
		ctx context.Context,
		accountClient internal.AccountIdWithAttachedClientttt,
		figies []internal.Figi,
		interval internal.TimeInterval,
	) ([]client.Operation, error)
}

type ProfitProvider struct {
	portfolioProvider    PortfolioProvider
	historyPriceProvider HistoryPriceBatchProvider
	operationsProvider   OperationsProvider
}

func NewProfitProvider(
	portfolioProvider PortfolioProvider,
	priceProvider HistoryPriceBatchProvider,
	operationsProvider OperationsProvider,
) *ProfitProvider {
	return &ProfitProvider{
		portfolioProvider:    portfolioProvider,
		historyPriceProvider: priceProvider,
		operationsProvider:   operationsProvider,
	}
}

func (p *ProfitProvider) Profits(
	ctx context.Context,
	accountClient internal.AccountIdWithAttachedClientttt,
	interval internal.TimeInterval,
	periods int,
) ([]internal.PeriodProfit, error) {
	operations, err := p.operationsProvider.Operations(
		ctx,
		accountClient,
		nil,
		internal.TimeInterval{
			End: interval.End,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get portfolio: %w", err)
	}

	pointInfos, figies := p.processOperations(operations, interval, periods+1)
	prices, err := p.historyPriceProvider.HistoryPriceBatch(ctx, accountClient, figies, interval, periods+1)
	if len(prices) != len(figies) {
		return nil, fmt.Errorf("got invalid prices for stocks")
	}
	for _, periodPrices := range prices {
		if len(periodPrices) != periods+1 {
			return nil, fmt.Errorf("got invalid prices for stocks")
		}
	}

	timeDiff := time.Duration(float64(interval.End.Sub(interval.Start)) / float64(periods))
	answerItems := make([]internal.PeriodProfit, 0, periods)
	for i := 0; i < periods; i++ {
		periodAnswerItem := internal.PeriodProfit{}
		previousPoint := pointInfos[i]
		currentPoint := pointInfos[i+1]
		previousBalance := previousPoint.balance
		currentBalance := currentPoint.balance
		for figi, income := range currentPoint.figiesIncome {
			figiPrices := prices[figi]

			currentMoneyInFigi := figiPrices[i+1].Price * internal.UnitsWithNano(currentPoint.figiesWithQuantity[figi])
			previousMoneyInFigi := figiPrices[i].Price * internal.UnitsWithNano(previousPoint.figiesWithQuantity[figi])

			figiIncome := income + currentMoneyInFigi - previousMoneyInFigi

			previousBalance += internal.Money(previousMoneyInFigi)
			currentBalance += internal.Money(currentMoneyInFigi)

			periodAnswerItem.FigiesProfit = append(periodAnswerItem.FigiesProfit, internal.FigiProfit{
				Figi:   figi,
				Profit: internal.Money(figiIncome),
			})
		}

		periodAnswerItem.AllProfit = currentBalance - previousBalance - currentPoint.moneyAdded
		periodAnswerItem.Interval = internal.TimeInterval{
			Start: interval.Start.Add(timeDiff * time.Duration(i)),
			End:   interval.Start.Add(timeDiff * time.Duration(i+1)),
		}

		answerItems = append(answerItems, periodAnswerItem)
	}

	return answerItems, nil
}

type pointInfo struct {
	figiesWithQuantity map[internal.Figi]internal.Quantity
	balance            internal.Money
	figiesIncome       map[internal.Figi]internal.UnitsWithNano
	moneyAdded         internal.Money
}

func (p *ProfitProvider) processOperations(operations []client.Operation, interval internal.TimeInterval, points int) ([]pointInfo, []internal.Figi) {
	pointInfos := make([]pointInfo, points)
	diff := time.Duration(float64(interval.End.Sub(interval.Start)) / float64(points-1))
	periodIndex := 0
	figiesHashSet := make(map[string]struct{}, 256)
	balance := internal.UnitsWithNano(0)
	figiesWithQuantity := make(map[string]int, 128)
	intervalStart := interval.Start
	moneyAdded := internal.UnitsWithNano(0)
	for _, operation := range operations {
		balance += operation.Payment
		if operation.Time.Before(intervalStart) {
			if operation.OperationType == int(internal.OperationTypeBuy) && operation.Figi != "" {
				figiesWithQuantity[operation.Figi] += operation.Quantity
			}

			if operation.OperationType == int(internal.OperationTypeSell) && operation.Figi != "" {
				figiesWithQuantity[operation.Figi] -= operation.Quantity
			}

			if operation.Figi != "" && periodIndex != 0 {
				pointInfos[periodIndex].figiesIncome[internal.Figi(operation.Figi)] += operation.Payment
			}

			if p.isMoneyAdding(operation) {
				moneyAdded += operation.Payment
			}
		} else {
			pointFigies := make(map[internal.Figi]internal.Quantity, len(figiesWithQuantity))
			for figi, qty := range figiesWithQuantity {
				if qty < 1 {
					if qty < 0 {
						fmt.Println("quantity is negative, fix it")
					}
					continue
				}

				figiesHashSet[figi] = struct{}{}
				pointFigies[internal.Figi(figi)] = internal.Quantity(qty)
			}

			pointInfos[periodIndex].balance = internal.Money(balance)
			pointInfos[periodIndex].moneyAdded = internal.Money(moneyAdded)
			moneyAdded = 0
			pointInfos[periodIndex].figiesWithQuantity = pointFigies

			periodIndex++
			if periodIndex >= points {
				break
			}
			pointInfos[periodIndex].figiesIncome = make(map[internal.Figi]internal.UnitsWithNano, int(float32(len(figiesHashSet))*1.25))

			intervalStart = intervalStart.Add(diff)
		}
	}

	uniqueFigies := make([]internal.Figi, 0, len(figiesHashSet))
	for figi := range figiesHashSet {
		uniqueFigies = append(uniqueFigies, internal.Figi(figi))
	}

	return pointInfos, uniqueFigies
}

func (p *ProfitProvider) isMoneyAdding(operation client.Operation) (yes bool) {
	yes = true
	switch operation.OperationType {
	case int(internal.OperationTypeInputSwift):
	case int(internal.OperationTypeInputAcquiring):
	case int(internal.OperationTypeInpMulti):
	case int(internal.OperationTypeInput):
	case int(internal.OperationTypeInputSecurities):
	case int(internal.OperationTypeOutput):
	case int(internal.OperationTypeOutputAcquiring):
	case int(internal.OperationTypeOutputSwift):
	case int(internal.OperationTypeOutMulti):
	default:
		yes = false
	}

	return
}
