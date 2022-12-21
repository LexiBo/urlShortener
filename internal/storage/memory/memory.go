package memory

import (
	"context"
	"github.com/inspectorvitya/shortenerURL/internal/storage"
	"sync"
)

type StorageInMemory struct {
	mu        sync.RWMutex
	data      map[string]string
	shortURLS map[string]struct{}
}

func New() *StorageInMemory {
	return &StorageInMemory{
		data:      make(map[string]string),
		shortURLS: make(map[string]struct{}),
	}
}

func (s *StorageInMemory) CreateShortURL(_ context.Context, fullURL, shortURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.data[shortURL]
	if ok {
		return storage.ErrAlreadyExists
	}
	s.data[shortURL] = fullURL
	s.shortURLS[fullURL] = struct{}{}
	return nil
}

func (s *StorageInMemory) GetFullURL(_ context.Context, shortURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fullURL, ok := s.data[shortURL]
	if !ok {
		return "", storage.ErrNotFound
	}
	return fullURL, nil
}

func (s *StorageInMemory) CheckExistFullURL(_ context.Context, fullURL string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.shortURLS[fullURL]
	return ok, nil
}
func (s *StorageInMemory) CheckExistShortURL(_ context.Context, shortURL string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.data[shortURL]
	return ok, nil
}
