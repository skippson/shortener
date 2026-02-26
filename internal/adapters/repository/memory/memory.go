package memory

import (
	"context"
	"shortener/internal/domain"
	"sync"
)

type MemoryRepository struct {
	mu             sync.RWMutex
	originalRepo   map[string]string
	shorteneddRepo map[string]string
}

func NewRepository() *MemoryRepository {
	return &MemoryRepository{
		mu:             sync.RWMutex{},
		originalRepo:   make(map[string]string),
		shorteneddRepo: make(map[string]string),
	}
}

func (r *MemoryRepository) Save(_ context.Context, original, shortened string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.originalRepo[original]; ok {
		return domain.ErrAlreadyExist
	}

	if _, ok := r.shorteneddRepo[shortened]; ok {
		return domain.ErrAlreadyExist
	}

	r.originalRepo[original] = shortened
	r.shorteneddRepo[shortened] = original

	return nil
}

func (r *MemoryRepository) GetByShortened(_ context.Context, shortened string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	original, ok := r.shorteneddRepo[shortened]
	if !ok {
		return "", domain.ErrNotFound
	}

	return original, nil
}

func (r *MemoryRepository) GetByOriginal(_ context.Context, original string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	shortened, ok := r.originalRepo[original]
	if !ok {
		return "", domain.ErrNotFound
	}

	return shortened, nil
}

func (r *MemoryRepository) Close() {}
