package rateupdaterworker

import (
	"context"
	"time"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/logger"
)

type usecase interface {
	UpdateCurrency(context.Context) error
}

type config interface {
	GetFrequencyRateUpdateSec() int
}

type RatesUpdaterWorker struct {
	usecase usecase
	cfg     config
}

func New(usecase usecase, cfg config) *RatesUpdaterWorker {
	return &RatesUpdaterWorker{
		usecase: usecase,
		cfg:     cfg,
	}
}

func (w RatesUpdaterWorker) Run(ctx context.Context) {
	err := w.usecase.UpdateCurrency(ctx)
	if err != nil {
		logger.Errorf("can not update currency: %v", err)
	}

	ticker := time.NewTicker(time.Duration(w.cfg.GetFrequencyRateUpdateSec()) * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			select {
			case <-ctx.Done():
				return
			default:
				err := w.usecase.UpdateCurrency(ctx)
				if err != nil {
					logger.Errorf("can not update currency: %v", err)
				}
			}
		}
	}
}
