package service

import (
	"context"
	"errors"
	"net/url"
	"shortLink/internal/model"
	"shortLink/internal/storage"
)

type Service struct {
	storage storage.Storage
}

func NewService(storage storage.Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) ShortenURL(ctx context.Context, originalURL string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	u, err := url.ParseRequestURI(originalURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return "", model.ErrInvalidUrl
	}

	for i := 0; i < 3; i++ {
		shortU, err := Generate()
		if err != nil {
			return "", err
		}
		shorterU, err := s.storage.Save(ctx, originalURL, shortU)
		if err == nil {
			return shorterU, nil
		}

		if errors.Is(err, model.ErrShortUrlCollision) {
			continue
		}
		return "", err
	}
	return "", errors.New("failed to generate shorten URL")
}

func (s *Service) GetOriginalURL(ctx context.Context, shortU string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	return s.storage.Get(ctx, shortU)
}
