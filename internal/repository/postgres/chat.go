package postgres

import "github.com/alonsoF100/chat-api/internal/models"

func (r Repository) CreateChat(title string) (*models.Chat, error) {
	const op = "postgres/chat.go/CreateChat"

	return nil, nil
}
func (r Repository) DeleteChat(ChatID int) error {
	const op = "postgres/chat.go/DeleteChat"

	return nil
}

func (r Repository) GetChat(ChatID int) (*models.Chat, error) {
	const op = "postgres/chat.go/GetChat"

	return nil, nil
}
