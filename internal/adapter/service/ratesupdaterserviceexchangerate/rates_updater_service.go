package ratesupdaterserviceexchangerate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
)

var errUnknownCurrencyCode = errors.New("unknown code")

type RatesUpdaterService struct{}

func New() *RatesUpdaterService {
	return &RatesUpdaterService{}
}

type RatesExch struct {
	Date  string
	Rates map[string]float64
}

func (s RatesUpdaterService) Get(ctx context.Context, base string, codes []string) ([]entity.Rate, error) {
	url := "https://api.exchangerate.host/latest"
	url = fmt.Sprintf("%s?base=%s&symbols=%s", url, base, strings.Join(codes, ","))

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

	var ratesExch RatesExch

	err = json.Unmarshal(body, &ratesExch)
	if err != nil {
		return nil, errors.Wrap(err, "RatesUpdaterService.Get")
	}

	timeParse, err := time.Parse("2006-01-02", ratesExch.Date)
	if err != nil {
		return nil, errors.Wrap(err, "RatesUpdaterService.Get")
	}

	time := entity.NewDateTimeFromTime(timeParse)

	rates := make([]entity.Rate, 0, len(codes))
	rates = append(rates, entity.NewRate(base, entity.NewDecimal(1, 0), time))

	for _, code := range codes {
		ratio, ok := ratesExch.Rates[code]
		if !ok {
			return nil, errUnknownCurrencyCode
		}

		rates = append(rates, entity.NewRate(code, entity.NewDecimalFromFloat(ratio), time))
	}

	return rates, nil
}
