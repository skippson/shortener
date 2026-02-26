package repository

import (
	"context"
)

type Repository interface {
	Save(ctx context.Context, original, shortened string) error
	GetByShortened(ctx context.Context, shortened string) (string, error)
	GetByOriginal(ctx context.Context, origin string) (string, error)
	Close()
}
