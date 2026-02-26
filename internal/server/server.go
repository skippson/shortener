package server

import (
	"context"
	"time"

	httphandlers "shortener/internal/controllers/http_handlers"
	"shortener/internal/controllers/http_handlers/middleware"
	"shortener/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	app *fiber.App
	log logger.Logger
}

func NewServer(h *httphandlers.ApiHandlers, log logger.Logger) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  10 * time.Second,
	})

	addHealthCheck(app)

	api := app.Group("/api")

	mw := middleware.NewMiddleware(log)
	h.MapApiRoutes(api, mw)

	return &Server{
		app: app,
		log: log,
	}
}

func addHealthCheck(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})
}

func (s *Server) Run(ctx context.Context, address string) error {
	errCh := make(chan error, 1)

	go func() {
		s.log.Info("server started",
			logger.Field{Key: "addres", Value: address},
		)

		if err := s.app.Listen(address); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		s.log.Info("shutting down server")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		return s.app.ShutdownWithContext(shutdownCtx)
	}
}
