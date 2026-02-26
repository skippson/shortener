package httphandlers

import (
	"shortener/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type Middleware interface {
	SetRequestID() fiber.Handler
}

type SuccessResponse[T any] struct {
	Data T `json:"data"`
}

type ErrorResponse struct {
	Msg    string `json:"message"`
	Status int    `json:"status"`
}

func writeSuccess[T any](c *fiber.Ctx, status int, data T) error {
	return c.Status(status).JSON(SuccessResponse[T]{
		Data: data,
	})
}

func writeError(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(ErrorResponse{
		Status: status,
		Msg:    msg,
	})
}

func getLogger(c *fiber.Ctx) logger.Logger {
	return c.Locals("logger").(logger.Logger)
}
