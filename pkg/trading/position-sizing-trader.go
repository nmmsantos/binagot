package trading

import (
	"time"

	"github.com/nmmsantos/binagot/pkg/account"
	"github.com/nmmsantos/binagot/pkg/bot"
)

type PositionSizingTrader struct {
	buyThreshold  float64
	sellThreshold float64
}

func NewPositionSizingTrader(buyThreshold float64, sellThreshold float64) *PositionSizingTrader {
	return &PositionSizingTrader{
		buyThreshold:  buyThreshold,
		sellThreshold: sellThreshold,
	}
}

func (t *PositionSizingTrader) Trade(acc *account.Accout) []bot.OperationRequest {
	now := time.Now().UTC()
	reqs := []bot.OperationRequest{}

	for _, s := range acc.Symbols {
		if s.RefPrice != 0 {
			percentage := s.Price/s.RefPrice - 1

			if percentage >= t.sellThreshold {
				reqs = append(reqs, bot.OperationRequest{
					Operation: bot.Sell,
					Symbol:    s,
					Buckets:   1,
				})
			} else if percentage <= -t.buyThreshold {
				reqs = append(reqs, bot.OperationRequest{
					Operation: bot.Buy,
					Symbol:    s,
					Buckets:   1,
				})
			} else {
				continue
			}
		}

		s.RefPrice = s.Price
		s.RefDate = now
	}

	return reqs
}
