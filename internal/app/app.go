package app

import (
	"context"
	"errors"
	"net/url"

	"github.com/inspectorvitya/shortenerURL/internal/storage"
	"github.com/inspectorvitya/shortenerURL/pkg/strgen"
)

type App struct {
	Storage storage.Storage
}

var (
	ErrInvalidURL = errors.New("invalid url")
	ErrExistURL   = errors.New("url exist")
)

func New(db storage.Storage) *App {
	return &App{Storage: db}
}

func (app *App) CreateShortURL(ctx context.Context, fullURL string) (string, error) {
	var token string
	var err error
	if !ValidURL(fullURL) {
		return "", ErrInvalidURL
	}
	exit, err := app.Storage.CheckExistFullURL(ctx, fullURL)
	if err != nil {
		return "", err
	}
	if exit {
		return "", ErrExistURL
	}
	for {
		token, err = strgen.GenerateRandomString(10)
		if err != nil {
			return "", err
		}
		exit, err := app.Storage.CheckExistShortURL(ctx, token)
		if err != nil {
			return "", err
		}
		if !exit {
			break
		}
	}
	err = app.Storage.CreateShortURL(ctx, fullURL, token)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (app *App) GetFullURL(ctx context.Context, shortURL string) (string, error) {
	fullURL, err := app.Storage.GetFullURL(ctx, shortURL)
	if err != nil {
		return "", err
	}
	return fullURL, nil
}

func ValidURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	return err == nil && u.Scheme != "" && u.Host != ""
}
