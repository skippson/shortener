package usecase

import (
	"context"
	"errors"
	"shortener/internal/domain"
)

type Repository interface {
	Save(ctx context.Context, original, short string) error
	GetByShortened(ctx context.Context, short string) (string, error)
	GetByOriginal(ctx context.Context, original string) (string, error)
}

type Generator interface {
	Generate() (string, error)
}

type Validator interface {
	ValidateURL(url string) (string, bool)
	ValidateShortened(shortened string) bool
}

type UsecaseOptions struct {
	Repository  Repository
	Generator   Generator
	Validator   Validator
	MaxAttempts int
	Protection  bool
}

type Usecase struct {
	repo        Repository
	gen         Generator
	validator   Validator
	maxAttempts int
	protec      bool
}

func NewUsecase(options UsecaseOptions) (*Usecase, error) {
	if options.MaxAttempts <= 0 {
		return nil, errors.New("maxAttempts must be positive")
	}
	return &Usecase{
		repo:        options.Repository,
		gen:         options.Generator,
		validator:   options.Validator,
		maxAttempts: options.MaxAttempts,
		protec:      options.Protection,
	}, nil
}

func (uc *Usecase) CreateShortened(ctx context.Context, url string) (string, error) {
	if uc.protec {
		ok := false
		url, ok = uc.validator.ValidateURL(url)
		if !ok {
			return "", domain.ErrInvalidURL
		}
	}

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

	return "", errors.New("maxAttempts exceeded")
}

func (uc *Usecase) GetOriginalByShortened(ctx context.Context, shortened string) (string, error) {
	if uc.protec {
		if !uc.validator.ValidateShortened(shortened) {
			return "", domain.ErrInvalidShortened
		}
	}

	origin, err := uc.repo.GetByShortened(ctx, shortened)
	if err != nil {
		if err == domain.ErrNotFound {
			return "", err
		}

		return "", err
	}

	return origin, nil
}
