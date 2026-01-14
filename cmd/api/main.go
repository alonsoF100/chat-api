package main

import (
	"log/slog"

	"github.com/alonsoF100/chat-api/internal/config"
	"github.com/alonsoF100/chat-api/internal/logger"
	"github.com/alonsoF100/chat-api/internal/repository/postgres"
	"github.com/alonsoF100/chat-api/internal/service"
	"github.com/alonsoF100/chat-api/internal/transport/http/handlers"
	"github.com/alonsoF100/chat-api/internal/transport/http/server"
	_ "github.com/alonsoF100/chat-api/migrations/postgres" // migrations
)

func main() {
	cfg := config.Load()

	logS := logger.Setup(cfg)

	pool, err := postgres.NewPool(cfg)
	if err != nil {
		slog.Error("Failed to create pool", "error", err)
	}
	defer pool.Close()
	slog.Info("Pool created successfully")

	dataBase := postgres.New(pool)

	service := service.New(dataBase)

	handlers := handlers.New(service)

	if err := server.New(cfg, handlers, logS).Start(); err != nil {
		slog.Error("Failed to start server",
			slog.String("error", err.Error()),
		)
	}
}
