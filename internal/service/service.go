package service

import "github.com/alonsoF100/chat-api/internal/models"

type Repository interface {
	CreateChat(title string) (*models.Chat, error)
	DeleteChat(ChatID int) error
	CreateMessage(ChatID int, text string) (*models.Message, error)
	GetMessages(ChatID int, limit int) ([]*models.Message, error)
	GetChat(ChatID int) (*models.Chat, error)
}

type Service struct {
	Repository Repository
}

func New(repository Repository) *Service {
	return &Service{
		Repository: repository,
	}
}
