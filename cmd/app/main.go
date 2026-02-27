package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"shortener/config"
	"shortener/internal/adapters/repository"
	"shortener/internal/adapters/repository/memory"
	"shortener/internal/adapters/repository/postgres"
	httphandlers "shortener/internal/controllers/http_handlers"
	"shortener/internal/generator"
	"shortener/internal/server"
	"shortener/internal/usecase"
	"shortener/pkg/logger"
	"syscall"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log, err := logger.New(cfg.Service.Name)
	if err != nil {
		panic(err)
	}

	log.Info("service starts working")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var db repository.Repository
	if cfg.Service.InMemory {
		db = memory.NewRepository()
	} else {
		db, err = postgres.NewRepository(ctx, cfg.Postgres)
		if err != nil {
			log.Error("database initialization error",
				logger.Field{Key: "error", Value: err})

			return
		}
		defer db.Close()
	}

	generator := generator.NewGenerator(cfg.Generator.Alphabet, cfg.Generator.Len)

	uc, err := usecase.NewUsecase(log, db, generator, cfg.Service.MaxGeneratorAttempts)
	if err != nil {
		log.Error("usecase initialization error",
			logger.Field{Key: "error", Value: err})

		return
	}

	apiControllers := httphandlers.NewHandlers(uc)

	srv := server.NewServer(apiControllers, log)

	if err = srv.Run(ctx, fmt.Sprintf("%s:%d", cfg.Service.Host, cfg.Service.Port)); err != nil {
		log.Error("server died",
			logger.Field{Key: "error", Value: err})

		return
	}

	log.Error("service successfully stopped")
}
