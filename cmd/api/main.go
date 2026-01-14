package main

import (
	"log/slog"

	"github.com/alonsoF100/chat-api/internal/config"
	"github.com/alonsoF100/chat-api/internal/logger"
	"github.com/alonsoF100/chat-api/internal/transport/http/handlers"
	"github.com/alonsoF100/chat-api/internal/transport/http/server"
)

func main() {
	cfg := config.Load()

	logS := logger.Setup(cfg)

	// TODO засетапить подключение к базе

	// TODO проинициализировать сервисный слой

	
	handlers := handlers.New(nil, nil) // TODO передать сервисный слой в http не забыть :)

	if err := server.New(cfg, handlers, logS).Start(); err != nil {
		slog.Error("Failed to start server",
			slog.String("error", err.Error()),
		)
	}
}
