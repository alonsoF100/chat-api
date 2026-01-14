package main

import (
	"log/slog"

	"github.com/alonsoF100/chat-api/internal/config"
	"github.com/alonsoF100/chat-api/internal/logger"
	"github.com/alonsoF100/chat-api/internal/service"
	"github.com/alonsoF100/chat-api/internal/transport/http/handlers"
	"github.com/alonsoF100/chat-api/internal/transport/http/server"
)

func main() {
	cfg := config.Load()

	logS := logger.Setup(cfg)

	// TODO засетапить подключение к базе

	service := service.New(nil) // TODO передать слой репо

	handlers := handlers.New(service)

	if err := server.New(cfg, handlers, logS).Start(); err != nil {
		slog.Error("Failed to start server",
			slog.String("error", err.Error()),
		)
	}
}
