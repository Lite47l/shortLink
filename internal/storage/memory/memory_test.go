package memory

import (
	"context"
	"errors"
	"shortLink/internal/model"
	"testing"
)

func TestMemory_Save(t *testing.T) {
	ctx := context.Background()
	mem := New()

	tests := []struct {
		name         string
		originalURL  string
		shortURL     string
		wantErr      error
		wantShortURL string
	}{
		{
			name:         "Успешное сохранение",
			originalURL:  "https://youtube.com",
			shortURL:     "Ab3_xY9zQw",
			wantErr:      nil,
			wantShortURL: "Ab3_xY9zQw",
		},
		{
			name:         "Идемпотентность - тот же оригинал возвращает ту же короткую ссылку",
			originalURL:  "https://youtube.com",
			shortURL:     "XyZ_123456",
			wantErr:      nil,
			wantShortURL: "Ab3_xY9zQw",
		},
		{
			name:         "Коллизия короткой ссылки",
			originalURL:  "https://google.com",
			shortURL:     "Ab3_xY9zQw",
			wantErr:      model.ErrShortUrlCollision,
			wantShortURL: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mem.Save(ctx, tt.originalURL, tt.shortURL)

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

			if result != tt.wantShortURL {
				t.Errorf("ожидался shortURL=%s, получен %s", tt.wantShortURL, result)
			}
		})
	}
}

func TestMemory_Get(t *testing.T) {
	ctx := context.Background()
	mem := New()

	mem.Save(ctx, "https://youtube.com", "Ab3_xY9zQw")

	tests := []struct {
		name     string
		shortURL string
		wantErr  error
		wantURL  string
	}{
		{
			name:     "Успешное получение",
			shortURL: "Ab3_xY9zQw",
			wantErr:  nil,
			wantURL:  "https://youtube.com",
		},
		{
			name:     "URL не найден",
			shortURL: "nonexistent",
			wantErr:  model.ErrNotFound,
			wantURL:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mem.Get(ctx, tt.shortURL)

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

			if result != tt.wantURL {
				t.Errorf("ожидался URL=%s, получен %s", tt.wantURL, result)
			}
		})
	}
}

func TestMemory_ConcurrentAccess(t *testing.T) {
	ctx := context.Background()
	mem := New()

	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func(id int) {
			originalURL := "https://example.com/" + string(rune('A'+id))
			shortURL := "Short" + string(rune('0'+id))

			mem.Save(ctx, originalURL, shortURL)
			mem.Get(ctx, shortURL)

			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}
