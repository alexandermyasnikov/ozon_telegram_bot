package currencymemorystorage

import (
	"errors"
	"sync"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
)

var errNoCurrencyFound = errors.New("can not find currency")

type CurrencyMemoryStorage struct {
	mutex sync.RWMutex
	rates map[string]entity.Rate
}

func New() *CurrencyMemoryStorage {
	return &CurrencyMemoryStorage{
		mutex: sync.RWMutex{},
		rates: map[string]entity.Rate{},
	}
}

func (s *CurrencyMemoryStorage) Get(currency string) (entity.Rate, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	rate, ok := s.rates[currency]
	if !ok {
		return entity.Rate{}, errNoCurrencyFound
	}

	return rate, nil
}

func (s *CurrencyMemoryStorage) GetAll() ([]entity.Rate, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	rates := make([]entity.Rate, 0, len(s.rates))
	for _, rate := range s.rates {
		rates = append(rates, rate)
	}

	return rates, nil
}

func (s *CurrencyMemoryStorage) Update(rate entity.Rate) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.rates[rate.GetCode()] = rate

	return nil
}
