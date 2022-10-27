package ratesupdaterservicecbr_test

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/service/ratesupdaterservicecbr"
)

func TestRatesUpdaterServiceCBR(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skip test with http request")
	}

	service := ratesupdaterservicecbr.New()
	ctx := context.Background()

	base := "RUB"
	codes := []string{"USD", "CNY", "EUR"}

	rates, err := service.Get(ctx, base, codes)

	assert.NoError(t, err)
	assert.Len(t, rates, 4)

	codesExpected := append([]string{}, codes...)
	codesExpected = append(codesExpected, base)

	codesActual := make([]string, 0, len(rates))
	for _, rate := range rates {
		codesActual = append(codesActual, rate.GetCode())
	}

	assert.ElementsMatch(t, codesExpected, codesActual)
}

func TestRatesUpdaterServiceCBR_TimeUpdate(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("skip test with http request")
	}

	service := ratesupdaterservicecbr.New()
	ctx := context.Background()

	base := "RUB"
	codes := []string{"USD", "CNY", "EUR"}

	rates, err := service.Get(ctx, base, codes)

	assert.NoError(t, err)
	assert.Len(t, rates, 4)

	time.Sleep(10 * time.Millisecond)

	rates2, err := service.Get(ctx, base, codes)
	assert.NoError(t, err)

	sort.Slice(rates, func(i, j int) bool {
		return rates[i].GetCode() < rates[j].GetCode()
	})

	sort.Slice(rates2, func(i, j int) bool {
		return rates2[i].GetCode() < rates2[j].GetCode()
	})

	for i := 0; i < len(rates); i++ {
		assert.Greater(t, rates2[i].GetTime(), rates[i].GetTime())
	}
}
