package ratesupdaterservicecbr

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
)

var errUnknownCurrencyCode = errors.New("unknown code")

type RatesUpdaterService struct{}

func New() *RatesUpdaterService {
	return &RatesUpdaterService{}
}

type RatesCBR struct {
	Date      string
	Timestamp int64
	Base      string
	Rates     map[string]float64
}

func (s RatesUpdaterService) Get(ctx context.Context, base string, codes []string) ([]entity.Rate, error) {
	url := "https://www.cbr-xml-daily.ru/latest.js"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "RatesUpdaterService.Get")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "RatesUpdaterService.Get")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "RatesUpdaterService.Get")
	}

	var ratesCBR RatesCBR

	err = json.Unmarshal(body, &ratesCBR)
	if err != nil {
		return nil, errors.Wrap(err, "RatesUpdaterService.Get")
	}

	time := entity.NewDateTimeFromTime(time.Unix(ratesCBR.Timestamp, 0))

	rates := make([]entity.Rate, 0, len(codes))
	rates = append(rates, entity.NewRate(base, entity.NewDecimal(1, 0), time))

	for _, code := range codes {
		ratio, ok := ratesCBR.Rates[code]
		if !ok {
			return nil, errUnknownCurrencyCode
		}

		rates = append(rates, entity.NewRate(code, entity.NewDecimalFromFloat(ratio), time))
	}

	return rates, nil
}
