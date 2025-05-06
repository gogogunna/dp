package intervals_provider

import (
	"context"
	"dp/internal"
	"dp/internal/client"
	"dp/pkg/nullable"
	"errors"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"golang.org/x/sync/errgroup"
	"math"
	"time"
)

var (
	source = investapi.GetCandlesRequest_CandleSource(investapi.GetCandlesRequest_CandleSource_value["CANDLE_SOURCE_INCLUDE_WEEKEND"])
)

const (
	goroutinesMax = 3
)

type intervalWithLimit struct {
	interval       time.Duration
	limit          int
	candleInterval investapi.CandleInterval
}

var availableIntervals = []intervalWithLimit{
	{1 * time.Minute, 2400, investapi.CandleInterval_CANDLE_INTERVAL_1_MIN},    // CANDLE_INTERVAL_1_MIN
	{2 * time.Minute, 1200, investapi.CandleInterval_CANDLE_INTERVAL_2_MIN},    // CANDLE_INTERVAL_2_MIN
	{3 * time.Minute, 750, investapi.CandleInterval_CANDLE_INTERVAL_3_MIN},     // CANDLE_INTERVAL_3_MIN
	{5 * time.Minute, 2400, investapi.CandleInterval_CANDLE_INTERVAL_5_MIN},    // CANDLE_INTERVAL_5_MIN
	{10 * time.Minute, 1200, investapi.CandleInterval_CANDLE_INTERVAL_10_MIN},  // CANDLE_INTERVAL_10_MIN
	{15 * time.Minute, 2400, investapi.CandleInterval_CANDLE_INTERVAL_15_MIN},  // CANDLE_INTERVAL_15_MIN
	{30 * time.Minute, 1200, investapi.CandleInterval_CANDLE_INTERVAL_30_MIN},  // CANDLE_INTERVAL_30_MIN
	{1 * time.Hour, 2400, investapi.CandleInterval_CANDLE_INTERVAL_HOUR},       // CANDLE_INTERVAL_HOUR
	{2 * time.Hour, 2400, investapi.CandleInterval_CANDLE_INTERVAL_2_HOUR},     // CANDLE_INTERVAL_2_HOUR
	{4 * time.Hour, 700, investapi.CandleInterval_CANDLE_INTERVAL_4_HOUR},      // CANDLE_INTERVAL_4_HOUR
	{24 * time.Hour, 2400, investapi.CandleInterval_CANDLE_INTERVAL_DAY},       // CANDLE_INTERVAL_DAY
	{7 * 24 * time.Hour, 300, investapi.CandleInterval_CANDLE_INTERVAL_WEEK},   // CANDLE_INTERVAL_WEEK
	{30 * 24 * time.Hour, 120, investapi.CandleInterval_CANDLE_INTERVAL_MONTH}, // CANDLE_INTERVAL_MONTH
}

type InstrumentHistoryPriceProvider struct{}

func NewInstrumentHistoryPriceProvider() *InstrumentHistoryPriceProvider {
	return &InstrumentHistoryPriceProvider{}
}

type candleDTO struct {
	time  time.Time
	price internal.UnitsWithNano
}

func (s *InstrumentHistoryPriceProvider) HistoryPrice(
	_ context.Context,
	accountClient internal.AccountIdWithAttachedClientttt,
	figi internal.Figi,
	interval internal.TimeInterval,
	points int,
) ([]client.InstrumentPrice, error) {
	if points == 0 {
		return nil, errors.New("points can't be 0")
	}
	marketClient := accountClient.Client.NewMarketDataServiceClient()
	diff := interval.End.Sub(interval.Start)
	queriesQuantity := math.MaxInt32
	for _, availableInterval := range availableIntervals {
		pointsQuantity := min(int(diff/availableInterval.interval+1), availableInterval.limit)
		if pointsQuantity < 2 {
			break
		}
		queriesQuantity = min(queriesQuantity, (points+pointsQuantity-1)/pointsQuantity)

		if queriesQuantity == 1 {
			break
		}
	}

	if queriesQuantity == math.MaxInt64 {
		return nil, errors.New("too ")
	}

	var maxPointsQuantity int
	index := nullable.Nullable[int]{}
	for i, availableInterval := range availableIntervals {
		pointsQuantity := min(int(diff/availableInterval.interval)+1, availableInterval.limit)

		tempQueriesQuantity := (points + pointsQuantity - 1) / pointsQuantity
		if tempQueriesQuantity == queriesQuantity {
			if maxPointsQuantity < pointsQuantity {
				maxPointsQuantity = pointsQuantity
				index.SetValue(i)
			}
		}
	}

	if index.IsNil() {
		return nil, errors.New("something went wrong while choosing interval")
	}

	stepInterval := availableIntervals[index.Value()].candleInterval
	limit := availableIntervals[index.Value()].limit
	candlesResp := make([][]*investapi.HistoricCandle, queriesQuantity)
	eg := errgroup.Group{}
	eg.SetLimit(goroutinesMax)
	for i := range queriesQuantity {
		eg.Go(func() error {
			resp, err := marketClient.GetCandles(
				string(figi),
				stepInterval,
				interval.Start, interval.End,
				source,
				int32(points-i*limit),
			)

			if err != nil {
				return err
			}

			candlesResp[i] = resp.GetCandles()

			return nil
		},
		)
	}

	mapped := mapCandles(candlesResp)

	if len(mapped) == 0 {
		return nil, nil
	}

	dynamicIndex := 0
	diffStep := diff / time.Duration(points)
	answer := make([]client.InstrumentPrice, 0, points)
	for i := range points {
		pointTime := interval.Start.Add(diffStep * time.Duration(i))
		if dynamicIndex >= len(mapped) {
			dynamicIndex = len(mapped) - 1
		}

		if mapped[dynamicIndex].time.Before(pointTime) {
			for ; mapped[dynamicIndex].time.Before(pointTime) && dynamicIndex > -1; dynamicIndex-- {
			}

			if dynamicIndex != len(mapped)-1 {
				dynamicIndex++
			}

			answer = append(answer, client.InstrumentPrice{
				Price:     mapped[dynamicIndex].price,
				RealTime:  mapped[dynamicIndex].time,
				PointTime: pointTime,
			})
		} else {
			for ; mapped[dynamicIndex].time.After(pointTime) && dynamicIndex < len(mapped); dynamicIndex++ {
			}

			if dynamicIndex != 0 {
				dynamicIndex--
			}

			answer = append(answer, client.InstrumentPrice{
				Price:     mapped[dynamicIndex].price,
				RealTime:  mapped[dynamicIndex].time,
				PointTime: pointTime,
			})
		}

	}

	return answer, nil
}

func (s *InstrumentHistoryPriceProvider) HistoryPric(
	ctx context.Context,
	accountClient internal.AccountIdWithAttachedClientttt,
	figi internal.Figi,
	interval internal.TimeInterval,
	points int,
) ([]client.InstrumentPrice, error) {
	if points <= 0 {
		return nil, errors.New("points must be positive")
	}

	marketClient := accountClient.Client.NewMarketDataServiceClient()
	duration := interval.End.Sub(interval.Start)

	// 1. Выбираем оптимальный интервал
	var bestInterval intervalWithLimit
	bestFit := math.MaxInt32

	for _, interval := range availableIntervals {
		totalPoints := int(duration/interval.interval) + 1
		if totalPoints < 2 {
			continue
		}

		// Выбираем интервал, который дает наиболее близкое количество точек
		if diff := abs(totalPoints - points); diff < bestFit {
			bestFit = diff
			bestInterval = interval
		}
	}

	if bestFit == math.MaxInt64 {
		return nil, errors.New("no suitable interval found for given time range")
	}

	// 2. Рассчитываем временные диапазоны для запросов с учетом лимита
	totalPossiblePoints := int(duration/bestInterval.interval) + 1
	queriesNeeded := (totalPossiblePoints + bestInterval.limit - 1) / bestInterval.limit
	pointsPerQuery := (totalPossiblePoints + queriesNeeded - 1) / queriesNeeded

	// 3. Параллельно получаем данные
	candlesResp := make([][]*investapi.HistoricCandle, queriesNeeded)
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(goroutinesMax)

	for i := 0; i < queriesNeeded; i++ {
		i := i
		eg.Go(func() error {
			// Рассчитываем временной диапазон для этого запроса
			pointsToGet := pointsPerQuery
			if i == queriesNeeded-1 {
				pointsToGet = totalPossiblePoints - i*pointsPerQuery
			}

			from := interval.Start.Add(time.Duration(i*pointsPerQuery) * bestInterval.interval)
			to := from.Add(time.Duration(pointsToGet) * bestInterval.interval)
			if to.After(interval.End) {
				to = interval.End
			}

			// Ограничиваем количество запрашиваемых точек лимитом
			limit := min(pointsToGet, bestInterval.limit)
			resp, err := marketClient.GetCandles(
				string(figi),
				bestInterval.candleInterval,
				from, to,
				source,
				int32(limit),
			)
			if err != nil {
				return err
			}

			candlesResp[i] = resp.GetCandles()
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// 4. Обрабатываем результаты
	allCandles := mapCandles(candlesResp)
	if len(allCandles) == 0 {
		return nil, nil
	}

	// 5. Создаем равномерно распределенные точки
	result := make([]client.InstrumentPrice, 0, points)
	step := duration / time.Duration(points-1)
	candleIndex := 0

	for i := 0; i < points; i++ {
		targetTime := interval.Start.Add(step * time.Duration(i))

		// Находим ближайшую свечу
		for candleIndex < len(allCandles)-1 &&
			allCandles[candleIndex+1].time.Before(targetTime) {
			candleIndex++
		}

		// Если есть следующая свеча и она ближе
		if candleIndex < len(allCandles)-1 &&
			targetTime.Sub(allCandles[candleIndex].time) >
				allCandles[candleIndex+1].time.Sub(targetTime) {
			candleIndex++
		}

		result = append(result, client.InstrumentPrice{
			Price:     allCandles[candleIndex].price,
			RealTime:  allCandles[candleIndex].time,
			PointTime: targetTime,
		})
	}

	return result, nil
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func mapCandles(candles [][]*investapi.HistoricCandle) []candleDTO {
	length := 0
	for _, candlesItem := range candles {
		length += len(candlesItem)
	}

	mapped := make([]candleDTO, 0, length)

	for _, candlesItem := range candles {
		for _, candle := range candlesItem {
			mapped = append(mapped, candleDTO{
				time:  candle.GetTime().AsTime(),
				price: internal.MapUnitsWithNano(candle.GetHigh()),
			})
		}
	}

	return mapped
}
