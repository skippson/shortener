package httphandlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *ApiHandlers) MapApiRoutes(router fiber.Router, mw Middleware) {
	router.Use(mw.SetRequestID())

	router.Post("/create_shortener", h.CreateShortened())
	router.Get("get_original/:shortener", h.GetOriginalal())
}
