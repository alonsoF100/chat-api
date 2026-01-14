package service

import "github.com/alonsoF100/chat-api/internal/models"

type Repository interface {
	CreateChat(title string) (*models.Chat, error)
	DeleteChat(chatID int) error
	CreateMessage(chatID int, text string) (*models.Message, error)
	GetMessages(chatID int, limit int) ([]*models.Message, error)
	GetChat(chatID int) (*models.Chat, error)
}

type Service struct {
	Repository Repository
}

func New(repository Repository) *Service {
	return &Service{
		Repository: repository,
	}
}
