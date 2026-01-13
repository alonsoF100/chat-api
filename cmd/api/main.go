package main

import (
	"fmt"

	"github.com/alonsoF100/chat-api/internal/config"
)

func main() {
	// TODO Инициализация зависимостей

	cfg := config.Load()
	fmt.Println(cfg)
}
