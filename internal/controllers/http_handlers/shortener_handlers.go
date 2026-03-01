package httphandlers

import (
	"context"
	"errors"

	"shortener/internal/domain"
	"shortener/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type Usecase interface {
	CreateShortened(ctx context.Context, url string) (string, error)
	GetOriginalByShortened(ctx context.Context, shortened string) (string, error)
}

type ApiHandlers struct {
	uc Usecase
}

func NewHandlers(uc Usecase) *ApiHandlers {
	return &ApiHandlers{
		uc: uc,
	}
}

type createShortenerParams struct {
	URL string `json:"url"`
}

type createShortenerResponse struct {
	Shortened string `json:"shortened"`
}

func (h *ApiHandlers) CreateShortened() fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := createShortenerParams{}
		if err := c.BodyParser(&req); err != nil {
			return writeError(c, fiber.StatusBadRequest, "invalid json")
		}

		shortened, err := h.uc.CreateShortened(c.Context(), req.URL)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidURL) {
				return writeError(c, fiber.StatusBadRequest, "invalid url")
			}

			getLogger(c).Error("create shortened failed",
				logger.Field{Key: "url", Value: req.URL},
				logger.Field{Key: "error", Value: err})

			return writeError(c, fiber.StatusInternalServerError, "internal error")
		}

		return writeSuccess(c, fiber.StatusOK, createShortenerResponse{Shortened: shortened})
	}
}

type getOriginalResponse struct {
	Original string `json:"original"`
}

func (h *ApiHandlers) GetOriginalal() fiber.Handler {
	return func(c *fiber.Ctx) error {
		shortened := c.Params("shortened")

		original, err := h.uc.GetOriginalByShortened(c.Context(), shortened)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidShortened) || errors.Is(err, domain.ErrNotFound) {
				return writeError(c, fiber.StatusNotFound, "not found")
			}

			getLogger(c).Error("get original failed",
				logger.Field{Key: "shortened", Value: shortened},
				logger.Field{Key: "error", Value: err})

			return writeError(c, fiber.StatusInternalServerError, "internal error")
		}

		return writeSuccess(c, fiber.StatusOK, getOriginalResponse{Original: original})
	}
}
