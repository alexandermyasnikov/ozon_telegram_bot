package currencymemorystorage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/currencymemorystorage"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
)

func TestCurrencyMemoryStorage(t *testing.T) {
	t.Parallel()

	dateTime1 := entity.NewDateTime(2022, 10, 9, 12, 00, 00)
	dateTime2 := entity.NewDateTime(2022, 10, 9, 12, 10, 00)

	storage := currencymemorystorage.New()

	_, err := storage.Get("RUB")
	assert.Error(t, err)

	err = storage.Update(entity.NewRate("RUB", entity.NewDecimal(1, 0), dateTime1))
	assert.NoError(t, err)

	_, err = storage.Get("EUR")
	assert.Error(t, err)
	assert.EqualError(t, err, "can not find currency")

	err = storage.Update(entity.NewRate("EUR", entity.NewDecimal(16327197, 9), dateTime1))
	assert.NoError(t, err)

	rate, err := storage.Get("EUR")
	assert.NoError(t, err)
	assert.Equal(t, entity.NewRate("EUR", entity.NewDecimal(16327197, 9), dateTime1), rate)

	err = storage.Update(entity.NewRate("EUR", entity.NewDecimal(16327121, 9), dateTime2))
	assert.NoError(t, err)

	rate, err = storage.Get("EUR")
	assert.NoError(t, err)
	assert.Equal(t, entity.NewRate("EUR", entity.NewDecimal(16327121, 9), dateTime2), rate)

	rate, err = storage.Get("RUB")
	assert.NoError(t, err)
	assert.Equal(t, entity.NewRate("RUB", entity.NewDecimal(1, 0), dateTime1), rate)

	rates, err := storage.GetAll()
	assert.NoError(t, err)
	assert.ElementsMatch(t, []entity.Rate{
		entity.NewRate("RUB", entity.NewDecimal(1, 0), dateTime1),
		entity.NewRate("EUR", entity.NewDecimal(16327121, 9), dateTime2),
	}, rates)
}
