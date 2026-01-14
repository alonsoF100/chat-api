package service

import (
	"log/slog"

	"github.com/alonsoF100/chat-api/internal/apperrors"
	"github.com/alonsoF100/chat-api/internal/models"
)

func (s Service) CreateMessage(text string, chatID int) (*models.Message, error) {
	const op = "service/CreateMessage"

	chat, err := s.Repository.GetChat(chatID)
	if chat == nil {
		slog.Debug("Failed to find chat",
			slog.String("op", op),
			slog.Int("chatID", chatID),
		)
		return nil, apperrors.ErrChatNotFound
	}
	if err != nil {
		slog.Error("Failed to find chat",
			slog.String("op", op),
			slog.Int("chatID", chatID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	message, err := s.Repository.CreateMessage(chatID, text)
	if err != nil {
		slog.Error("Failed to create message",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	slog.Info("Message created successfully",
		slog.String("op", op),
		slog.Int("messageID", message.ID),
	)
	return message, nil
}

func (s Service) GetMessages(chatID int, limit int) (*models.ChatWithMessages, error) {
	const op = "service/GetMessages"

	if limit > 100 {
		limit = 100
	}
	if limit < 0 {
		limit = 20
	}

	chat, err := s.Repository.GetChat(chatID)
	if chat == nil {
		slog.Debug("Failed to find chat",
			slog.String("op", op),
			slog.Int("chatID", chatID),
		)
		return nil, apperrors.ErrChatNotFound
	}
	if err != nil {
		slog.Error("Failed to find chat",
			slog.String("op", op),
			slog.Int("chatID", chatID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	messages, err := s.Repository.GetMessages(chatID, limit)
	if err != nil {
		slog.Error("Failed to get mesages from db",
			slog.String("op", op),
			slog.Int("chatID", chatID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	slog.Info("Messages founded successfully",
		slog.String("op", op),
		slog.Any("chat", chat),
		slog.Int("messagesCount", len(messages)),
	)
	return &models.ChatWithMessages{Chat: chat, Messages: messages}, nil
}
