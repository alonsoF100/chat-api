package main

import (
	"github.com/alonsoF100/chat-api/internal/config"
	"github.com/alonsoF100/chat-api/internal/logger"
)

func main() {
	// TODO Инициализация зависимостей

	cfg := config.Load()

	logger.Setup(cfg)
}
