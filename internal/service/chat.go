package service

import (
	"log/slog"
	"strings"

	"github.com/alonsoF100/chat-api/internal/apperrors"
	"github.com/alonsoF100/chat-api/internal/models"
)

func (s Service) CreateChat(title string) (*models.Chat, error) {
	const op = "service/CreateChat"

	validTitle := strings.TrimSpace(title)

	chat, err := s.Repository.CreateChat(validTitle)
	if err != nil {
		slog.Error("Failed to create chat",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	slog.Info("Chat created successfully",
		slog.String("op", op),
		slog.Int("chatID", chat.ID),
	)
	return chat, nil
}

func (s Service) DeleteChat(chatID int) error {
	const op = "service/DeleteChat"

	chat, err := s.Repository.GetChat(chatID)
	if chat == nil {
		slog.Debug("Failed to find chat",
			slog.String("op", op),
			slog.Int("chatID", chatID),
		)
		return apperrors.ErrChatNotFound
	}
	if err != nil {
		slog.Error("Failed to find chat",
			slog.String("op", op),
			slog.Int("chatID", chatID),
			slog.String("error", err.Error()),
		)
		return err
	}

	err = s.Repository.DeleteChat(chatID)
	if err != nil {
		slog.Error("Failed to delete chat",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return err
	}

	slog.Info("Chat deleted successfully",
		slog.String("op", op),
		slog.Int("chatID", chatID),
	)
	return nil
}
