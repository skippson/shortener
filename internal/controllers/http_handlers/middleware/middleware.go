package middleware

import (
	"shortener/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Middleware struct {
	log logger.Logger
}

func NewMiddleware(log logger.Logger) *Middleware {
	return &Middleware{
		log: log,
	}
}

func (mw *Middleware) SetRequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		logWithReq := mw.log.With(logger.Field{Key: "request_id", Value: uuid.New().String()})

		c.Locals("logger", logWithReq)

		return c.Next()
	}
}
