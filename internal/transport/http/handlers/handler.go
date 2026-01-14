package handlers

import (
	"github.com/alonsoF100/chat-api/internal/models"
	"github.com/go-playground/validator/v10"
)

type ChatService interface {
	CreateChat(title string) (*models.Chat, error)
	DeleteChat(chatID int) error
}

type MessageService interface {
	CreateMessage(text string, chatID int) (*models.Message, error)
	GetMessages(chatID int, limit int) ([]*models.Message, error)
}

type Handler struct {
	ChatService    ChatService
	MessageService MessageService
	Validator      *validator.Validate
}

func New(chatService ChatService, messageService MessageService) *Handler {
	return &Handler{
		ChatService:    chatService,
		MessageService: messageService,
		Validator:      validator.New(),
	}
}
