package handlers

import (
	"github.com/alonsoF100/chat-api/internal/models"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	CreateChat(title string) (*models.Chat, error)
	DeleteChat(chatID int) error
	CreateMessage(text string, chatID int) (*models.Message, error)
	GetMessages(chatID int, limit int) (*models.ChatWithMessages, error)
}

type Handler struct {
	Service   Service
	Validator *validator.Validate
}

func New(service Service) *Handler {
	return &Handler{
		Service:   service,
		Validator: validator.New(),
	}
}
