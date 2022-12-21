package storage

import (
	"context"
	"errors"
)

var (
	ErrAlreadyExists = errors.New("url already exists")
	ErrNotFound      = errors.New("url not found")
)

type Storage interface {
	CreateShortURL(ctx context.Context, fullURL, shortURL string) error
	GetFullURL(ctx context.Context, shortURL string) (string, error)
	CheckExistFullURL(ctx context.Context, fullURL string) (bool, error)
	CheckExistShortURL(ctx context.Context, shortURL string) (bool, error)
}
