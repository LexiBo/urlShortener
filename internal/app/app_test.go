package app

import (
	"context"
	"testing"

	"github.com/inspectorvitya/shortenerURL/internal/storage"
	"github.com/stretchr/testify/require"
)

type testStorage struct {
	existURL     map[string]string
	existFullURL map[string]struct{}
}

func (s *testStorage) CreateShortURL(_ context.Context, fullURL, shortURL string) error {
	_, ok := s.existURL[shortURL]
	if ok {
		return storage.ErrAlreadyExists
	}
	s.existURL[shortURL] = fullURL
	s.existFullURL[fullURL] = struct{}{}
	return nil
}
func (s *testStorage) GetFullURL(_ context.Context, shortURL string) (string, error) {
	fullURL, ok := s.existURL[shortURL]
	if !ok {
		return "", storage.ErrNotFound
	}
	return fullURL, nil
}
func (s *testStorage) CheckExistFullURL(_ context.Context, fullURL string) (bool, error) {
	_, ok := s.existFullURL[fullURL]
	return ok, nil
}
func (s *testStorage) CheckExistShortURL(_ context.Context, shortURL string) (bool, error) {
	_, ok := s.existURL[shortURL]
	return ok, nil
}

func TestApp_CreateShortUrl(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		db := &testStorage{
			existURL:     make(map[string]string),
			existFullURL: make(map[string]struct{}),
		}
		shorter := &App{db}
		shortURL, err := shorter.CreateShortURL(context.Background(), "https://github.com/")
		require.NoError(t, err)
		require.NotEmpty(t, shortURL)
	})
	t.Run("invalid url", func(t *testing.T) {
		db := &testStorage{
			existURL:     make(map[string]string),
			existFullURL: make(map[string]struct{}),
		}
		shorter := &App{db}
		shortURL, err := shorter.CreateShortURL(context.Background(), "gasd21da")
		require.ErrorIs(t, ErrInvalidURL, err)
		require.Empty(t, shortURL)
	})
	t.Run("url exist", func(t *testing.T) {
		db := &testStorage{
			existURL:     make(map[string]string),
			existFullURL: make(map[string]struct{}),
		}
		db.existFullURL["https://github.com/"] = struct{}{}
		shorter := &App{db}
		shortURL, err := shorter.CreateShortURL(context.Background(), "https://github.com/")
		require.ErrorIs(t, ErrExistURL, err)
		require.Empty(t, shortURL)
	})
}
