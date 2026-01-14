package postgres

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/alonsoF100/chat-api/internal/models"
)

func (r Repository) CreateMessage(chatID int, text string) (*models.Message, error) {
	const op = "postgres/message.go/CreateMessage"

	const query = `
	INSERT INTO messages (chat_id, text) 
	VALUES ($1, $2) 
	RETURNING id, chat_id, text, created_at
	`
	var message models.Message
	err := r.pool.QueryRow(context.Background(), query, chatID, text).Scan(
		&message.ID,
		&message.ChatID,
		&message.Text,
		&message.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &message, nil
}

func (r Repository) GetMessages(chatID int, limit int) ([]*models.Message, error) {
	const op = "postgres/message.go/GetMessages"

	query, args, err := squirrel.Select("id", "chat_id", "text", "created_at").
		From("messages").
		Where(squirrel.Eq{"chat_id": chatID}).
		OrderBy("created_at DESC").
		Limit(uint64(limit)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := r.pool.Query(context.Background(), query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		var message models.Message
		err := rows.Scan(
			&message.ID,
			&message.ChatID,
			&message.Text,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		messages = append(messages, &message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return messages, nil
}
