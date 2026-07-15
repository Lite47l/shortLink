package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"shortLink/internal/model"
	"shortLink/internal/service"
)

type Handler struct {
	service *service.Service
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/shorten", h.Shorten)
	mux.HandleFunc("GET /{id}", h.GetOriginal)
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	if err := r.Context().Err(); err != nil {
		http.Error(w, "request cancelled", http.StatusServiceUnavailable)
		return
	}

	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	shortUrl, err := h.service.ShortenURL(r.Context(), req.URL)
	if err != nil {
		if errors.Is(err, model.ErrInvalidUrl) {
			http.Error(w, "invalid url", http.StatusBadRequest)
			return
		}

		if errors.Is(err, r.Context().Err()) {
			http.Error(w, "request canceled", http.StatusServiceUnavailable)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	log.Printf("Short URL created: %s", shortUrl)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shortenResponse{ShortURL: shortUrl})
}

func (h *Handler) GetOriginal(w http.ResponseWriter, r *http.Request) {
	if err := r.Context().Err(); err != nil {
		http.Error(w, "request cancelled", http.StatusServiceUnavailable)
		return
	}

	shortID := r.PathValue("id")
	if shortID == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	originalUrl, err := h.service.GetOriginalURL(r.Context(), shortID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		if errors.Is(err, r.Context().Err()) {
			http.Error(w, "request canceled", http.StatusServiceUnavailable)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	log.Printf("Original URL found: %s", originalUrl)
	http.Redirect(w, r, originalUrl, http.StatusFound)
}
