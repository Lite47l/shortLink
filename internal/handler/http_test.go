package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"shortLink/internal/model"
	"shortLink/internal/service"
	"testing"
)

type mockStorageForHandler struct {
	data map[string]string
}

func (m *mockStorageForHandler) Save(ctx context.Context, originalURL, shortURL string) (string, error) {
	if m.data == nil {
		m.data = make(map[string]string)
	}
	m.data[shortURL] = originalURL
	return shortURL, nil
}

func (m *mockStorageForHandler) Get(ctx context.Context, shortURL string) (string, error) {
	original, ok := m.data[shortURL]
	if !ok {
		return "", model.ErrNotFound
	}
	return original, nil
}

func (m *mockStorageForHandler) Close() error {
	return nil
}

func TestHandler_Shorten(t *testing.T) {
	mockStore := &mockStorageForHandler{
		data: make(map[string]string),
	}
	svc := service.NewService(mockStore)
	h := NewHandler(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "Успешное создание",
			body:       `{"url":"https://youtube.com"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "Невалидный JSON",
			body:       `invalid json`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Пустой URL",
			body:       `{"url":""}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Невалидный URL",
			body:       `{"url":"not-a-url"}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ожидался статус %d, получен %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestHandler_GetOriginal(t *testing.T) {
	mockStore := &mockStorageForHandler{
		data: map[string]string{
			"Ab3_xY9zQw": "https://youtube.com",
		},
	}
	svc := service.NewService(mockStore)
	h := NewHandler(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantLoc    string
	}{
		{
			name:       "Успешный редирект",
			path:       "/Ab3_xY9zQw",
			wantStatus: http.StatusFound,
			wantLoc:    "https://youtube.com",
		},
		{
			name:       "URL не найден",
			path:       "/nonexistent",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			mux.ServeHTTP(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("ожидался статус %d, получен %d", tt.wantStatus, w.Code)
			}

			if tt.wantLoc != "" {
				loc := w.Header().Get("Location")
				if loc != tt.wantLoc {
					t.Errorf("ожидался Location=%s, получен %s", tt.wantLoc, loc)
				}
			}
		})
	}
}
