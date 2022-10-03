package productmemorystorage

import (
	"sync"
	"time"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/storage"
)

type ProductMemoryStorage struct {
	mutex    sync.RWMutex
	products map[int64]map[int64][]entity.Product
}

var _ storage.ProductStorageInterface = (*ProductMemoryStorage)(nil)

func NewProductStorage() *ProductMemoryStorage {
	return &ProductMemoryStorage{
		mutex:    sync.RWMutex{},
		products: make(map[int64]map[int64][]entity.Product),
	}
}

func (s *ProductMemoryStorage) Add(userID int64, product entity.Product) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.products[userID] == nil {
		s.products[userID] = make(map[int64][]entity.Product)
	}

	date := product.GetDate()
	day := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	s.products[userID][day.Unix()] = append(s.products[userID][day.Unix()], product)

	return nil
}

func (s *ProductMemoryStorage) GetAll(userID int64, date time.Time, days int) ([]entity.Product, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	products := make([]entity.Product, 0)

	productsUser, ok := s.products[userID]
	if !ok {
		return products, nil
	}

	day := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	for i := 0; i < days; i++ {
		productsDay, ok := productsUser[day.Unix()]
		if ok {
			// TODO урорядочить по дате, тут, при вставке или в service
			products = append(products, productsDay...)
		}

		day = day.AddDate(0, 0, -1)
	}

	return products, nil
}
