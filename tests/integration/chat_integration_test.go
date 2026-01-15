package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alonsoF100/chat-api/internal/config"
	"github.com/alonsoF100/chat-api/internal/repository/postgres"
	"github.com/alonsoF100/chat-api/internal/service"
	"github.com/alonsoF100/chat-api/internal/transport/http/dto"
	"github.com/alonsoF100/chat-api/internal/transport/http/handlers"
	"github.com/alonsoF100/chat-api/internal/transport/http/router"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatFullLifecycle_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:         8080,
			ReadTimeout:  time.Second * 5,
			WriteTimeout: time.Second * 10,
			IdleTimeout:  time.Second * 10,
		},
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "postgres",
			Name:     "chat-api",
			SSLMode:  "disable",
		},
		Logger: config.LoggerConfig{
			Level: "info",
			JSON:  false,
		},
		Migration: config.MigrationsConfig{
			Dir: "../../migrations/postgres",
		},
	}

	ctx := context.Background()

	pool, err := postgres.NewPool(cfg)
	require.NoError(t, err, "Failed to create database pool")
	defer pool.Close()

	err = pool.Ping(ctx)
	require.NoError(t, err, "Failed to ping database")

	database := postgres.New(pool)
	service := service.New(database)
	handler := handlers.New(service)

	router := router.New(handler).Setup()
	httpServer := httptest.NewServer(router)
	defer httpServer.Close()


	// Создаем чат через API
	chatData := dto.CreateChatRequest{Title: "Integration Test Chat"}
	chatBody, _ := json.Marshal(chatData)

	resp, err := http.Post(httpServer.URL+"/chats/", "application/json", bytes.NewBuffer(chatBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode, "Failed to create chat")

	var chatResp dto.ChatResponse
	err = json.NewDecoder(resp.Body).Decode(&chatResp)
	require.NoError(t, err)
	require.NotZero(t, chatResp.ID, "Chat ID should not be zero")

	chatID := chatResp.ID
	resp.Body.Close()

	// Добавляем первое сообщение
	message1Data := dto.CreateMessageRequest{Text: "First integration test message"}
	message1Body, _ := json.Marshal(message1Data)

	resp, err = http.Post(fmt.Sprintf("%s/chats/%d/messages/", httpServer.URL, chatID), "application/json", bytes.NewBuffer(message1Body))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode, "Failed to create first message")

	var message1Resp dto.MessageResponse
	err = json.NewDecoder(resp.Body).Decode(&message1Resp)
	require.NoError(t, err)
	require.NotZero(t, message1Resp.ID, "Message ID should not be zero")

	message1ID := message1Resp.ID
	resp.Body.Close()

	// Добавляем второе сообщение
	message2Data := dto.CreateMessageRequest{Text: "Second integration test message"}
	message2Body, _ := json.Marshal(message2Data)

	resp, err = http.Post(fmt.Sprintf("%s/chats/%d/messages/", httpServer.URL, chatID), "application/json", bytes.NewBuffer(message2Body))
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode, "Failed to create second message")

	var message2Resp dto.MessageResponse
	err = json.NewDecoder(resp.Body).Decode(&message2Resp)
	require.NoError(t, err)
	require.NotZero(t, message2Resp.ID, "Message ID should not be zero")

	message2ID := message2Resp.ID
	resp.Body.Close()

	// Получаем сообщения из чата
	resp, err = http.Get(fmt.Sprintf("%s/chats/%d", httpServer.URL, chatID))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, "Failed to get messages")

	var messagesResp dto.MessagesResponse
	err = json.NewDecoder(resp.Body).Decode(&messagesResp)
	require.NoError(t, err)

	// Проверяем информацию о чате
	assert.Equal(t, chatID, messagesResp.Chat.ID, "Chat ID should match")
	assert.Equal(t, "Integration Test Chat", messagesResp.Chat.Title, "Chat title should match")

	// Проверяем что есть хотя бы 2 сообщения
	assert.GreaterOrEqual(t, len(messagesResp.Messages), 2, "Should have at least 2 messages")

	// Ищем созданные сообщения
	foundMessage1 := false
	foundMessage2 := false
	for _, msg := range messagesResp.Messages {
		if msg.ID == message1ID && msg.Text == "First integration test message" {
			foundMessage1 = true
		}
		if msg.ID == message2ID && msg.Text == "Second integration test message" {
			foundMessage2 = true
		}
	}
	assert.True(t, foundMessage1, "First message should be in response")
	assert.True(t, foundMessage2, "Second message should be in response")

	resp.Body.Close()

	// Пробуем добавить сообщение в несуществующий чат
	nonExistentChatID := 99999
	invalidMessageData := dto.CreateMessageRequest{Text: "Should fail"}
	invalidMessageBody, _ := json.Marshal(invalidMessageData)

	resp, err = http.Post(fmt.Sprintf("%s/chats/%d/messages/", httpServer.URL, nonExistentChatID), "application/json", bytes.NewBuffer(invalidMessageBody))
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode, "Should return 404 for non-existent chat")
	resp.Body.Close()

	// Удаляем чат через API
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/chats/%d", httpServer.URL, chatID), nil)
	require.NoError(t, err)

	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode, "Failed to delete chat")
	resp.Body.Close()

	// Проверяем что сообщения удалились каскадно
	messages, err := database.GetMessages(chatID, 10)
	require.NoError(t, err)
	assert.Empty(t, messages, "Messages should be deleted when chat is deleted")

	// Проверяем что нельзя получить сообщения из удаленного чата
	resp, err = http.Get(fmt.Sprintf("%s/chats/%d", httpServer.URL, chatID))
	require.NoError(t, err)

	// Должен вернуть 404 чата уже нету
	require.Equal(t, http.StatusNotFound, resp.StatusCode, "Deleted chat should return 404")
	resp.Body.Close()
}
