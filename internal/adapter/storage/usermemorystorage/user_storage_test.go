package usermemorystorage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/usermemorystorage"
)

func TestUserMemoryStorage(t *testing.T) {
	t.Parallel()

	storage := usermemorystorage.New()

	_, err := storage.GetDefaultCurrency(1)
	assert.Error(t, err)
	assert.EqualError(t, err, "can not find user")

	err = storage.UpdateDefaultCurrency(1, "EUR")
	assert.NoError(t, err)

	currency, err := storage.GetDefaultCurrency(1)
	assert.NoError(t, err)
	assert.Equal(t, "EUR", currency)

	err = storage.UpdateDefaultCurrency(1, "RUB")
	assert.NoError(t, err)

	currency, err = storage.GetDefaultCurrency(1)
	assert.NoError(t, err)
	assert.Equal(t, "RUB", currency)
}
