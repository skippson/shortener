package usecase_test

import (
	"context"
	"errors"
	"testing"

	"shortener/internal/domain"
	"shortener/internal/usecase"
	"shortener/internal/usecase/internal/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateShortened(t *testing.T) {
	ctx := context.Background()
	maxAttempts := 2

	tests := []struct {
		name          string
		original      string
		wantShortened string
		setUpMocks    func(repo *mocks.MockRepository, gen *mocks.MockGenerator, validator *mocks.MockValidator)
		wantErr       assert.ErrorAssertionFunc
		protection    bool
	}{
		{
			name:          "ok",
			original:      "example",
			wantShortened: "ok",
			setUpMocks: func(repo *mocks.MockRepository, gen *mocks.MockGenerator, validator *mocks.MockValidator) {
				repo.EXPECT().GetByOriginal(ctx, "example").Return("", domain.ErrNotFound)
				gen.EXPECT().Generate().Return("ok", nil)
				repo.EXPECT().Save(ctx, "example", "ok").Return(nil)
			},
			wantErr:    assert.NoError,
			protection: false,
		},
		{
			name:          "protection",
			original:      "example",
			wantShortened: "",
			setUpMocks: func(repo *mocks.MockRepository, gen *mocks.MockGenerator, validator *mocks.MockValidator) {
				validator.EXPECT().ValidateURL("example").Return("", false)
			},
			wantErr:    assert.Error,
			protection: true,
		},
		{
			name:          "already exists",
			original:      "example",
			wantShortened: "exist",
			setUpMocks: func(repo *mocks.MockRepository, gen *mocks.MockGenerator, validator *mocks.MockValidator) {
				repo.EXPECT().GetByOriginal(ctx, "example").Return("exist", nil)
			},
			wantErr:    assert.NoError,
			protection: false,
		},
		{
			name:          "fisrt db error",
			original:      "example",
			wantShortened: "",
			setUpMocks: func(repo *mocks.MockRepository, gen *mocks.MockGenerator, validator *mocks.MockValidator) {
				repo.EXPECT().GetByOriginal(ctx, "example").Return("", errors.New("db error"))
			},
			wantErr:    assert.Error,
			protection: false,
		},
		{
			name:          "second db error",
			original:      "example",
			wantShortened: "",
			setUpMocks: func(repo *mocks.MockRepository, gen *mocks.MockGenerator, validator *mocks.MockValidator) {
				repo.EXPECT().GetByOriginal(ctx, "example").Return("", domain.ErrNotFound)
				gen.EXPECT().Generate().Return("ok", nil)
				repo.EXPECT().Save(ctx, "example", "ok").Return(errors.New("db error"))
			},
			wantErr:    assert.Error,
			protection: false,
		},
		{
			name:          "generator error",
			original:      "example",
			wantShortened: "",
			setUpMocks: func(repo *mocks.MockRepository, gen *mocks.MockGenerator, validator *mocks.MockValidator) {
				repo.EXPECT().GetByOriginal(ctx, "example").Return("", domain.ErrNotFound)
				gen.EXPECT().Generate().Return("", errors.New("gen error"))
			},
			wantErr:    assert.Error,
			protection: false,
		},
		{
			name:          "collision",
			original:      "example",
			wantShortened: "ok",
			setUpMocks: func(repo *mocks.MockRepository, gen *mocks.MockGenerator, validator *mocks.MockValidator) {
				repo.EXPECT().GetByOriginal(ctx, "example").Return("", domain.ErrNotFound)
				gen.EXPECT().Generate().Return("collision", nil)
				repo.EXPECT().Save(ctx, "example", "collision").Return(domain.ErrAlreadyExist)
				repo.EXPECT().GetByOriginal(ctx, "example").Return("", domain.ErrNotFound)
				gen.EXPECT().Generate().Return("ok", nil)
				repo.EXPECT().Save(ctx, "example", "ok").Return(nil)
			},
			wantErr:    assert.NoError,
			protection: false,
		},
		{
			name:          "max attempts exceeded",
			original:      "example",
			wantShortened: "",
			setUpMocks: func(repo *mocks.MockRepository, gen *mocks.MockGenerator, validator *mocks.MockValidator) {
				for i := 0; i < maxAttempts; i++ {
					repo.EXPECT().GetByOriginal(ctx, "example").Return("", domain.ErrNotFound)
					gen.EXPECT().Generate().Return("collision", nil)
					repo.EXPECT().Save(ctx, "example", "collision").Return(domain.ErrAlreadyExist)
				}
			},
			wantErr:    assert.Error,
			protection: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRepository(ctrl)
			gen := mocks.NewMockGenerator(ctrl)
			logger := mocks.NewMockLogger(ctrl)
			validator := mocks.NewMockValidator(ctrl)

			logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
			tt.setUpMocks(repo, gen, validator)

			uc, _ := usecase.NewUsecase(usecase.UsecaseOptions{
				Repository:  repo,
				Generator:   gen,
				Validator:   validator,
				Logger:      logger,
				MaxAttempts: maxAttempts,
				Protection:  tt.protection,
			})

			gotShortened, err := uc.CreateShortened(ctx, tt.original)
			tt.wantErr(t, err)
			assert.Equal(t, tt.wantShortened, gotShortened)
		})
	}
}

func TestGetShortenedByOriginal(t *testing.T) {
	ctx := context.Background()
	maxAttempts := 1

	tests := []struct {
		name       string
		shortened  string
		mockRes    string
		wantValue  string
		mockErr    error
		wantErr    assert.ErrorAssertionFunc
		protection bool
	}{
		{
			name:       "ok",
			shortened:  "ok",
			mockRes:    "example",
			wantValue:  "example",
			mockErr:    nil,
			wantErr:    assert.NoError,
			protection: false,
		},
		{
			name:       "protection",
			shortened:  "enemy",
			mockRes:    "",
			wantValue:  "",
			mockErr:    domain.ErrInvalidShortened,
			wantErr:    assert.Error,
			protection: true,
		},
		{
			name:      "not found",
			shortened: "ok?",
			mockRes:   "",
			wantValue: "",
			mockErr:   domain.ErrNotFound,
			wantErr:   assert.Error,
		},
		{
			name:       "repo error",
			shortened:  "ok",
			mockRes:    "",
			wantValue:  "",
			mockErr:    errors.New("db error"),
			wantErr:    assert.Error,
			protection: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockRepository(ctrl)
			gen := mocks.NewMockGenerator(ctrl)
			logger := mocks.NewMockLogger(ctrl)
			validator := mocks.NewMockValidator(ctrl)

			logger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
			if tt.protection {
				validator.EXPECT().ValidateShortened(tt.shortened).Return(false)
			} else {
				repo.EXPECT().GetByShortened(ctx, tt.shortened).Return(tt.mockRes, tt.mockErr)
			}

			uc, _ := usecase.NewUsecase(usecase.UsecaseOptions{
				Repository:  repo,
				Generator:   gen,
				Validator:   validator,
				Logger:      logger,
				MaxAttempts: maxAttempts,
				Protection:  tt.protection,
			})

			gotOriginal, err := uc.GetShortenedByOriginal(ctx, tt.shortened)

			tt.wantErr(t, err)
			assert.Equal(t, tt.wantValue, gotOriginal)
		})
	}
}
