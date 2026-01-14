package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/alonsoF100/chat-api/internal/models"
)

func (r Repository) CreateChat(title string) (*models.Chat, error) {
	const op = "postgres/chat.go/CreateChat"

	const query = `
	INSERT INTO chats (title) 
	VALUES ($1) 
	RETURNING id, title, created_at
	`
	var chat models.Chat

	err := r.pool.QueryRow(context.Background(), query, title).Scan(
		&chat.ID,
		&chat.Title,
		&chat.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	slog.Debug("Chat created successfully",
		slog.String("op", op),
		slog.Int("id", chat.ID),
		slog.String("title", chat.Title),
		slog.Time("created_at", chat.CreatedAt),
	)
	return &chat, nil
}

func (r Repository) DeleteChat(chatID int) error {
	const op = "postgres/chat.go/DeleteChat"

	const query = `
	DELETE FROM chats 
	WHERE id = $1
	`
	res, err := r.pool.Exec(context.Background(), query, chatID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("%s: chat with id %d not found", op, chatID)
	}

	slog.Debug("Chat deleted successfully",
		slog.String("op", op),
		slog.Int("id", chatID),
	)
	return nil
}

func (r Repository) GetChat(chatID int) (*models.Chat, error) {
	const op = "postgres/chat.go/GetChat"

	const query = `
	SELECT id, title, created_at 
	FROM chats 
	WHERE id = $1
	`

	var chat models.Chat
	err := r.pool.QueryRow(context.Background(), query, chatID).Scan(
		&chat.ID,
		&chat.Title,
		&chat.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	slog.Debug("Chat founded successfully",
		slog.String("op", op),
		slog.Int("id", chat.ID),
		slog.String("title", chat.Title),
		slog.Time("created_at", chat.CreatedAt),
	)
	return &chat, nil
}
