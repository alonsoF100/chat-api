package dto

import (
	"time"

	"github.com/alonsoF100/chat-api/internal/models"
)


type ErrorResponse struct {
	Error     string    `json:"error"`
	Timestamp time.Time `json:"time_stamp"`
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Error:     err.Error(),
		Timestamp: time.Now(),
	}
}

type ChatResponse struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

func NewChatResponse(chat *models.Chat) ChatResponse {
	return ChatResponse{
		ID:        chat.ID,
		Title:     chat.Title,
		CreatedAt: chat.CreatedAt,
	}
}

type MessageResponse struct {
	ID        int       `json:"id"`
	ChatID    int       `json:"chat_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

func NewMessageResponse(message *models.Message) MessageResponse {
	return MessageResponse{
		ID:        message.ID,
		ChatID:    message.ChatID,
		Text:      message.Text,
		CreatedAt: message.CreatedAt,
	}
}

type MessagesResponse struct {
	Messages []*MessageResponse `json:"messages"`
}

func NewMessagesResponse(messages []*models.Message) MessagesResponse {
	responseMessages := MessagesResponse{
		Messages: make([]*MessageResponse, 0, len(messages)),
	}

	for _, message := range messages {
		temp := &MessageResponse{
			ID:        message.ID,
			ChatID:    message.ChatID,
			Text:      message.Text,
			CreatedAt: message.CreatedAt,
		}
		responseMessages.Messages = append(responseMessages.Messages, temp)
	}
	return responseMessages
}
