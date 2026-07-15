package memory

import (
	"context"
	"shortLink/internal/model"
	"sync"
)

type Memory struct {
	mu                 sync.RWMutex
	originalURLToShort map[string]string
	shortURLToOriginal map[string]string
}

func New() *Memory {
	return &Memory{
		originalURLToShort: make(map[string]string),
		shortURLToOriginal: make(map[string]string),
	}
}

func (m *Memory) Save(ctx context.Context, originalURL, shortURL string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if existingShortU, ok := m.originalURLToShort[originalURL]; ok {
		return existingShortU, nil
	}

	if existingOrigU, ok := m.shortURLToOriginal[shortURL]; ok && existingOrigU != originalURL {
		return "", model.ErrShortUrlCollision
	}

	m.originalURLToShort[originalURL] = shortURL
	m.shortURLToOriginal[shortURL] = originalURL
	return shortURL, nil
}

func (m *Memory) Get(ctx context.Context, shortURL string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	original, ok := m.shortURLToOriginal[shortURL]
	if !ok {
		return "", model.ErrNotFound
	}
	return original, nil
}

func (m *Memory) Close() error {
	return nil
}
