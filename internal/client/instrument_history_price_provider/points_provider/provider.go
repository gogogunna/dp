package points_provider

import (
	"context"
	"dp/internal"
	"dp/internal/client"
	"errors"
	"fmt"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"golang.org/x/sync/errgroup"
	"time"
)

const (
	goroutinesMax = 3
)

type InstrumentPointsPriceProvider = int

var (
	source = investapi.GetCandlesRequest_CandleSource(investapi.GetCandlesRequest_CandleSource_value["CANDLE_SOURCE_INCLUDE_WEEKEND"])
)

type InstrumentHistoryPriceProvider struct{}

func NewInstrumentHistoryPriceProvider() *InstrumentHistoryPriceProvider {
	return &InstrumentHistoryPriceProvider{}
}

type FigiWithPoint struct {
	Figi  internal.Figi
	Point time.Time
}

func (s *InstrumentHistoryPriceProvider) HistoryPrice(
	_ context.Context,
	accountClient internal.AccountIdWithAttachedClientttt,
	figiesWithPoints []FigiWithPoint,
) ([]client.InstrumentPrice, error) {
	marketClient := accountClient.Client.NewMarketDataServiceClient()

	eg := errgroup.Group{}
	eg.SetLimit(goroutinesMax)
	answer := make([]client.InstrumentPrice, len(figiesWithPoints))
	for i, figiWithPoint := range figiesWithPoints {
		eg.Go(func() error {
			resp, err := marketClient.GetCandles(
				string(figiWithPoint.Figi),
				investapi.CandleInterval_CANDLE_INTERVAL_4_HOUR,
				figiWithPoint.Point, figiWithPoint.Point.Add(time.Hour*24),
				source,
				100,
			)
			if err != nil {
				return fmt.Errorf("failed to get price for stock: %w", err)
			}

			if len(resp.GetCandles()) == 0 {
				return errors.New("failed to get price for stock")
			}

			candle := resp.GetCandles()[0]
			instrumentPrice := client.InstrumentPrice{
				Price:     internal.MapUnitsWithNano(candle.GetHigh()),
				RealTime:  candle.GetTime().AsTime(),
				PointTime: figiWithPoint.Point,
			}

			answer[i] = instrumentPrice

			return nil
		})
	}

	err := eg.Wait()

	return answer, err
}
