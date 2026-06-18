package memory

import (
	"sync"

	"github.com/go-thin/go-auth/pkg/models"
	"github.com/go-thin/go-auth/pkg/storage"
)

// InMemoryStorage is a thread-safe in-memory implementation of storage.Storage.
// Intended for testing and development only.
type InMemoryStorage struct {
	mu    sync.RWMutex
	users map[string]models.User // keyed by username
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{users: make(map[string]models.User)}
}

func (s *InMemoryStorage) CreateUser(user models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[user.Username]; exists {
		return storage.ErrUserNotFound
	}
	s.users[user.Username] = user
	return nil
}

func (s *InMemoryStorage) GetUserByUsername(username string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[username]
	if !ok {
		return nil, storage.ErrUserNotFound
	}
	return &u, nil
}

func (s *InMemoryStorage) DeleteUser(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for username, user := range s.users {
		if user.ID == id {
			delete(s.users, username)
			return nil
		}
	}
	return nil
}
