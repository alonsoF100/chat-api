package postgres

import "github.com/alonsoF100/chat-api/internal/models"

func (r Repository) CreateMessage(ChatID int, text string) (*models.Message, error) {
	const op = "postgres/message.go/CreateMessage"

	return nil, nil
}

func (r Repository) GetMessages(ChatID int, limit int) ([]*models.Message, error) {
	const op = "postgres/message.go/GetMessages"

	return nil, nil
}
