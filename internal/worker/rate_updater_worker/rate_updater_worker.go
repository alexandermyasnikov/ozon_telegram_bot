package rateupdaterworker

import (
	"context"
	"log"
	"time"
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
					log.Println(err)
				}
			}
		}
	}
}
