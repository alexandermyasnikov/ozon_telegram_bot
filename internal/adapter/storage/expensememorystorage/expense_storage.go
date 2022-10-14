package expensememorystorage

import (
	"sync"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
)

type ExpenseMemoryStorage struct {
	mutex    sync.RWMutex
	expenses map[key][]entity.Expense
}

type key struct {
	UserID entity.UserID
	Date   int64
}

func New() *ExpenseMemoryStorage {
	return &ExpenseMemoryStorage{
		mutex:    sync.RWMutex{},
		expenses: make(map[key][]entity.Expense),
	}
}

func (s *ExpenseMemoryStorage) Create(userID entity.UserID, expense entity.Expense) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	k := key{
		UserID: userID,
		Date:   expense.Getdate().ToInt64(),
	}

	s.expenses[k] = append(s.expenses[k], expense)

	return nil
}

func (s *ExpenseMemoryStorage) Get(userID entity.UserID, date entity.Date, days int) ([]entity.Expense, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	expenses := make([]entity.Expense, 0)

	for ; days > 0; days-- {
		k := key{
			UserID: userID,
			Date:   date.ToInt64(),
		}

		expenses = append(expenses, s.expenses[k]...)

		date.AddDays(-1)
	}

	/*sort.SliceStable(expenses, func(i, j int) bool {
		return expenses[i].GetCategory() < expenses[j].GetCategory()
	})*/

	return expenses, nil
}
