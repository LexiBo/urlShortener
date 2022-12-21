package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/inspectorvitya/shortenerURL/internal/storage"
)

func (s *StorageInDB) CreateShortURL(ctx context.Context, fullURL, shortURL string) error {
	query := `INSERT INTO urls (full_url, short_url) VALUES ($1, $2)`

	_, err := s.db.ExecContext(ctx, query, fullURL, shortURL)
	if err != nil {
		return err
	}

	return nil
}
func (s *StorageInDB) GetFullURL(ctx context.Context, shortURL string) (string, error) {
	query := `select full_url from urls where short_url = $1`
	var fullURL string
	err := s.db.QueryRowxContext(ctx, query, shortURL).Scan(&fullURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrNotFound
		}
	}
	return fullURL, nil
}
func (s *StorageInDB) CheckExistFullURL(ctx context.Context, fullURL string) (bool, error) {
	query := `select exists(select FROM urls WHERE full_url = $1)`
	var exist bool
	err := s.db.QueryRowxContext(ctx, query, fullURL).Scan(&exist)

	if err != nil {
		return false, fmt.Errorf("delete banner storage: %w", err)
	}

	return exist, nil
}
func (s *StorageInDB) CheckExistShortURL(ctx context.Context, shortURL string) (bool, error) {
	query := `select exists(select FROM urls WHERE short_url = $1)`
	var exist bool
	err := s.db.QueryRowxContext(ctx, query, shortURL).Scan(&exist)

	if err != nil {
		return false, fmt.Errorf("delete banner storage: %w", err)
	}

	return exist, nil
}
