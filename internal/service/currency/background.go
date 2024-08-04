package currency

import (
	"context"
	"time"
)

func (s *Currency) RunBackgroudProcesses(ctx context.Context) {
	go s.priceCheckLoop(ctx)
}

func (s *Currency) priceCheckLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.priceCheckTicker.C:
			_, err := s.GetCurrentPrices(ctx)
			if err != nil {
				s.logger.Error("priceCheckLoop: failed to get current prices: " + err.Error())
				s.priceCheckTicker.Reset(1 * time.Minute)
				continue
			}

			s.priceCheckTicker.Reset(s.priceCheckInterval)
		}
	}
}
