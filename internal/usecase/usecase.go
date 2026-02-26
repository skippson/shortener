package usecase

import (
	"context"
	"errors"
	"shortener/internal/domain"
	"shortener/pkg/logger"
)

type Repository interface {
	Save(ctx context.Context, original, short string) error
	GetByShortened(ctx context.Context, short string) (string, error)
	GetByOriginal(ctx context.Context, original string) (string, error)
}

type Generator interface {
	Generate() (string, error)
}

type Usecase struct {
	repo        Repository
	gen         Generator
	log         logger.Logger
	maxAttempts int
}

func NewUsecase(logger logger.Logger, repo Repository, gen Generator, maxAttempts int) (*Usecase, error) {
	if maxAttempts <= 0 {
		return nil, errors.New("maxAttempts must be positive")
	}
	return &Usecase{
		repo:        repo,
		gen:         gen,
		log:         logger,
		maxAttempts: maxAttempts,
	}, nil
}

func (uc *Usecase) CreateShortened(ctx context.Context, url string) (string, error) {
	for range uc.maxAttempts {
		shortened, err := uc.repo.GetByOriginal(ctx, url)
		if err == nil {
			return shortened, nil
		}

		if !errors.Is(err, domain.ErrNotFound) {
			return "", err
		}

		shortened, err = uc.gen.Generate()
		if err != nil {
			uc.log.Error("cannot generate shortened url",
				logger.Field{Key: "origin", Value: url},
				logger.Field{Key: "error", Value: err})

			return "", err
		}

		err = uc.repo.Save(ctx, url, shortened)
		if err != nil {
			if errors.Is(err, domain.ErrAlreadyExist) {
				continue
			}

			return "", err
		}

		return shortened, nil
	}

	uc.log.Error("max attempts exceeded",
		logger.Field{Key: "origin", Value: url})

	return "", errors.New("maxAttempts exceeded")
}

func (uc *Usecase) GetShortenedByOriginal(ctx context.Context, shortened string) (string, error) {
	origin, err := uc.repo.GetByShortened(ctx, shortened)
	if err != nil {
		if err == domain.ErrNotFound {
			return "", err
		}

		uc.log.Error("cannot get origin url",
			logger.Field{Key: "shortened", Value: shortened},
			logger.Field{Key: "error", Value: err})

		return "", err
	}

	return origin, nil
}
