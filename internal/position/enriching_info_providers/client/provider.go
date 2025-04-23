package client

import (
	"context"
	"dp/internal"
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync"
)

const (
	goroutinesMax = 10
)

type PositionEnrichingInfoProvider struct{}

func NewClientPositionEnrichingInfoProvider() *PositionEnrichingInfoProvider {
	return &PositionEnrichingInfoProvider{}
}

func (s *PositionEnrichingInfoProvider) PositionEnrichingInfo(
	_ context.Context,
	accountClient internal.AccountIdWithAttachedClientttt,
	figies []internal.Figi,
) (map[internal.Figi]internal.PositionEnrichingInfo, error) {
	instrumentsClient := accountClient.Client.NewInstrumentsServiceClient()

	eg := errgroup.Group{}
	answer := make(map[internal.Figi]internal.PositionEnrichingInfo, len(figies))
	mu := sync.Mutex{}
	eg.SetLimit(goroutinesMax)
	for _, figi := range figies {
		eg.Go(func() error {
			resp, err := instrumentsClient.InstrumentByFigi(string(figi))
			if err != nil {
				return err
			}

			if resp.InstrumentResponse == nil {
				return nil
			}

			domain := mapEnrichingInfoToDomain(resp.InstrumentResponse)
			mu.Lock()
			answer[figi] = domain
			mu.Unlock()

			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		return nil, fmt.Errorf("failed to get positions enriching info: %w", err)
	}

	return answer, nil
}
