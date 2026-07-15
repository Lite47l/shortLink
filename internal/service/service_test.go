package service

import (
	"context"
	"errors"
	"shortLink/internal/model"
	"testing"
)

type mockStorage struct {
	savedShortURL string
	savedErr      error
	getOriginal   string
	getErr        error
}

func (m *mockStorage) Save(ctx context.Context, originalURL, shortURL string) (string, error) {
	return m.savedShortURL, m.savedErr
}

func (m *mockStorage) Get(ctx context.Context, shortURL string) (string, error) {
	return m.getOriginal, m.getErr
}

func (m *mockStorage) Close() error {
	return nil
}

func TestService_ShortenURL(t *testing.T) {
	tests := []struct {
		name        string
		originalURL string
		mock        *mockStorage
		wantErr     error
		wantURL     string
	}{
		{
			name:        "Успешное создание ссылки",
			originalURL: "https://youtube.com",
			mock: &mockStorage{
				savedShortURL: "Ab3_xY9zQw",
				savedErr:      nil,
			},
			wantErr: nil,
			wantURL: "Ab3_xY9zQw",
		},
		{
			name:        "Невалидный URL",
			originalURL: "not-a-url",
			mock:        &mockStorage{},
			wantErr:     model.ErrInvalidUrl,
		},
		{
			name:        "URL без схемы",
			originalURL: "youtube.com",
			mock:        &mockStorage{},
			wantErr:     model.ErrInvalidUrl,
		},
		{
			name:        "Ошибка хранилища",
			originalURL: "https://google.com",
			mock: &mockStorage{
				savedErr: errors.New("db error"),
			},
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(tt.mock)
			ctx := context.Background()

			shortURL, err := svc.ShortenURL(ctx, tt.originalURL)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ожидалась ошибка %v, но получена nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Errorf("ожидалась ошибка %v, получена %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("неожиданная ошибка: %v", err)
				return
			}

			if shortURL != tt.wantURL {
				t.Errorf("ожидался shortURL=%s, получен %s", tt.wantURL, shortURL)
			}
		})
	}
}

func TestService_GetOriginalURL(t *testing.T) {
	tests := []struct {
		name     string
		shortURL string
		mock     *mockStorage
		wantErr  error
		wantURL  string
	}{
		{
			name:     "Успешное получение URL",
			shortURL: "Ab3_xY9zQw",
			mock: &mockStorage{
				getOriginal: "https://youtube.com",
				getErr:      nil,
			},
			wantErr: nil,
			wantURL: "https://youtube.com",
		},
		{
			name:     "URL не найден",
			shortURL: "nonexistent",
			mock: &mockStorage{
				getErr: model.ErrNotFound,
			},
			wantErr: model.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NewService(tt.mock)
			ctx := context.Background()

			originalURL, err := svc.GetOriginalURL(ctx, tt.shortURL)

			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("ожидалась ошибка %v, но получена nil", tt.wantErr)
					return
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("ожидалась ошибка %v, получена %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Errorf("неожиданная ошибка: %v", err)
				return
			}

			if originalURL != tt.wantURL {
				t.Errorf("ожидался URL=%s, получен %s", tt.wantURL, originalURL)
			}
		})
	}
}
