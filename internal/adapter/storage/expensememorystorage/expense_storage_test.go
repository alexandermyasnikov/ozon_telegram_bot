package expensememorystorage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/storage/expensememorystorage"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
)

func TestExpenseMemoryStorage(t *testing.T) {
	t.Parallel()

	storage := expensememorystorage.New()

	userID := entity.UserID(1)

	expensesInit := []entity.Expense{
		entity.NewExpense("cat5", entity.NewDecimal(1, 0), entity.NewDate(2022, 10, 9)),
		entity.NewExpense("cat4", entity.NewDecimal(2, 0), entity.NewDate(2022, 10, 8)),
		entity.NewExpense("cat5", entity.NewDecimal(2, 0), entity.NewDate(2022, 10, 7)),
		entity.NewExpense("cat1", entity.NewDecimal(3, 0), entity.NewDate(2022, 10, 4)),
		entity.NewExpense("cat5", entity.NewDecimal(4, 0), entity.NewDate(2022, 10, 4)),
		entity.NewExpense("cat1", entity.NewDecimal(5, 0), entity.NewDate(2022, 10, 2)),
		entity.NewExpense("cat1", entity.NewDecimal(6, 0), entity.NewDate(2022, 9, 10)),
		entity.NewExpense("cat1", entity.NewDecimal(7, 0), entity.NewDate(2022, 9, 3)),
	}

	for _, experse := range expensesInit {
		err := storage.Create(userID, experse)
		assert.NoError(t, err)
	}

	expenses, err := storage.Get(userID, entity.NewDate(2022, 10, 8), 1)
	assert.NoError(t, err)
	assert.Equal(t, []entity.Expense{
		entity.NewExpense("cat4", entity.NewDecimal(2, 0), entity.NewDate(2022, 10, 8)),
	}, expenses)

	expenses, err = storage.Get(userID, entity.NewDate(2022, 10, 8), 2)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []entity.Expense{
		entity.NewExpense("cat4", entity.NewDecimal(2, 0), entity.NewDate(2022, 10, 8)),
		entity.NewExpense("cat5", entity.NewDecimal(2, 0), entity.NewDate(2022, 10, 7)),
	}, expenses)

	expenses, err = storage.Get(userID, entity.NewDate(2022, 10, 8), 7)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []entity.Expense{
		entity.NewExpense("cat4", entity.NewDecimal(2, 0), entity.NewDate(2022, 10, 8)),
		entity.NewExpense("cat5", entity.NewDecimal(2, 0), entity.NewDate(2022, 10, 7)),
		entity.NewExpense("cat1", entity.NewDecimal(3, 0), entity.NewDate(2022, 10, 4)),
		entity.NewExpense("cat5", entity.NewDecimal(4, 0), entity.NewDate(2022, 10, 4)),
		entity.NewExpense("cat1", entity.NewDecimal(5, 0), entity.NewDate(2022, 10, 2)),
	}, expenses)

	expenses, err = storage.Get(userID, entity.NewDate(2022, 10, 8), 31)
	assert.NoError(t, err)
	assert.ElementsMatch(t, []entity.Expense{
		entity.NewExpense("cat4", entity.NewDecimal(2, 0), entity.NewDate(2022, 10, 8)),
		entity.NewExpense("cat5", entity.NewDecimal(2, 0), entity.NewDate(2022, 10, 7)),
		entity.NewExpense("cat1", entity.NewDecimal(3, 0), entity.NewDate(2022, 10, 4)),
		entity.NewExpense("cat5", entity.NewDecimal(4, 0), entity.NewDate(2022, 10, 4)),
		entity.NewExpense("cat1", entity.NewDecimal(5, 0), entity.NewDate(2022, 10, 2)),
		entity.NewExpense("cat1", entity.NewDecimal(6, 0), entity.NewDate(2022, 9, 10)),
	}, expenses)
}
