package usermemorystorage

import (
	"errors"
	"sync"

	"gitlab.ozon.dev/myasnikov.alexander.s/telegram-bot/internal/entity"
)

var errNoUserFound = errors.New("can not find user")

type UserMemoryStorage struct {
	mutex sync.RWMutex
	users map[entity.UserID]*entity.User
}

func New() *UserMemoryStorage {
	return &UserMemoryStorage{
		mutex: sync.RWMutex{},
		users: make(map[entity.UserID]*entity.User),
	}
}

func (s *UserMemoryStorage) GetDefaultCurrency(userID entity.UserID) (string, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, ok := s.users[userID]
	if !ok {
		return "", errNoUserFound
	}

	return user.GetDefaultCurrency(), nil
}

func (s *UserMemoryStorage) UpdateDefaultCurrency(userID entity.UserID, currency string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.users[userID] == nil {
		user := entity.NewUser(userID)
		s.users[userID] = &user
	}

	s.users[userID].SetDefaultCurrency(currency)

	return nil
}
