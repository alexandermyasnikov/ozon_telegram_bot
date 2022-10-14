package ratesupdaterservicecbr_test

import (
	"context"
	"testing"

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
