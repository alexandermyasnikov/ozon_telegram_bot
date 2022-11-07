package currencycachestorage

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	lrucache "gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/adapter/cache/lru"
	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
)

type CurrencyStorage interface {
	Get(context.Context, string) (entity.Rate, error)
	GetAll(context.Context) ([]entity.Rate, error)
	Update(context.Context, entity.Rate) error
}

type Config interface {
	GetCurrencyCacheEnable() bool
	GetCurrencyCacheSize() int
	GetCurrencyCacheTTL() int
}

type CurrencyCacheStorage struct {
	storage CurrencyStorage
	cfg     Config
	cache   *lrucache.LRUCache
}

func New(storage CurrencyStorage, cfg Config) *CurrencyCacheStorage {
	var cache *lrucache.LRUCache
	if cfg.GetCurrencyCacheEnable() {
		cache = lrucache.NewLRUCache(&sync.RWMutex{}, cfg.GetCurrencyCacheSize())
	}

	return &CurrencyCacheStorage{
		storage: storage,
		cfg:     cfg,
		cache:   cache,
	}
}

func (s *CurrencyCacheStorage) Get(ctx context.Context, currency string) (entity.Rate, error) {
	if rate, ok := s.getFromCache(currency); ok {
		return rate, nil
	}

	rate, err := s.storage.Get(ctx, currency)
	if err != nil {
		return entity.Rate{}, errors.Wrap(err, "CurrencyCacheStorage.Get")
	}

	s.addToCache(rate)

	return rate, nil
}

func (s *CurrencyCacheStorage) GetAll(ctx context.Context) ([]entity.Rate, error) {
	rate, err := s.storage.GetAll(ctx)

	return rate, errors.Wrap(err, "CurrencyCacheStorage.GetAll")
}

func (s *CurrencyCacheStorage) Update(ctx context.Context, rate entity.Rate) error {
	s.addToCache(rate)

	err := s.storage.Update(ctx, rate)

	return errors.Wrap(err, "CurrencyCacheStorage.Update")
}

func (s *CurrencyCacheStorage) addToCache(rate entity.Rate) {
	if !s.cfg.GetCurrencyCacheEnable() {
		return
	}

	s.cache.Add(time.Now(), rate.GetCode(), rate, s.cfg.GetCurrencyCacheTTL())
}

func (s *CurrencyCacheStorage) getFromCache(currency string) (entity.Rate, bool) {
	if !s.cfg.GetCurrencyCacheEnable() {
		return entity.Rate{}, false
	}

	val, ok := s.cache.Get(time.Now(), currency)
	if !ok {
		return entity.Rate{}, false
	}

	rate, ok := val.(entity.Rate)
	if !ok {
		return entity.Rate{}, false
	}

	return rate, true
}
