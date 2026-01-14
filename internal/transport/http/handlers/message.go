package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/alonsoF100/chat-api/internal/apperrors"
	"github.com/alonsoF100/chat-api/internal/transport/http/dto"
	"github.com/go-chi/chi/v5"
)

/*
pattern: /chats/{id}/messages/
method: POST
info: JSON in request body + chatID from URL

succeed:

	-status code: 201 created
	-response body: JSON represented created message

failed:

	-status code: 400 bad request, 404 not found, 500 internal server error
	-response body: JSON with error message + timestamp
*/
func (h Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	const op = "handlers/CreateMessage"

	var req dto.CreateMessageRequest

	id := chi.URLParam(r, "id")
	chatID, err := strconv.Atoi(id)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, dto.NewErrorResponse(err))
		slog.Debug("Failed to convert string to int",
			slog.String("op", op),
			slog.String("string", id),
			slog.String("error", err.Error()),
		)
		return
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, dto.NewErrorResponse(err))
		slog.Debug("Failed to decode",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		WriteJSON(w, http.StatusBadRequest, dto.NewErrorResponse(err))
		slog.Debug("Failed to validate JSON",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return
	}

	message, err := h.Service.CreateMessage(req.Text, chatID)
	if err != nil {
		if errors.Is(err, apperrors.ErrChatNotFound) {
			WriteJSON(w, http.StatusNotFound, dto.NewErrorResponse(err))
			slog.Error("Failed to find chat",
				slog.Int("chatID", chatID),
				slog.String("error", err.Error()),
			)
			return
		}

		WriteJSON(w, http.StatusInternalServerError, dto.NewErrorResponse(err))
		slog.Error("Failed to create message",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return
	}

	WriteJSON(w, http.StatusCreated, dto.NewMessageResponse(message))
	slog.Info("Message created successfuly",
		slog.Int("chatID", message.ChatID),
		slog.Int("messageID", message.ID),
	)
}

/*
pattern: /chats/{id}?limit=N, N beetween 20 and 100
method: GET
info: chatID from URL + limit from query params

succeed:

	-status code: 200 ok
	-response body: JSON represented chat + []messages

failed:

	-status code: 400 bad request, 404 not found, 500 internal server error
	-response body: JSON with error message + timestamp
*/
func (h Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	const op = "handlers/GetMessages"

	id := chi.URLParam(r, "id")
	chatID, err := strconv.Atoi(id)
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, dto.NewErrorResponse(err))
		slog.Debug("Failed to conver string to int",
			slog.String("op", op),
			slog.String("string", id),
			slog.String("error", err.Error()),
		)
		return
	}

	var limit int
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		limit = 20
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			WriteJSON(w, http.StatusBadRequest, dto.NewErrorResponse(err))
			slog.Debug("Failed to convert string to int",
				slog.String("op", op),
				slog.String("string", limitStr),
				slog.String("error", err.Error()),
			)
			return
		}
	}

	chatWithMessages, err := h.Service.GetMessages(chatID, limit)
	if err != nil {
		if errors.Is(err, apperrors.ErrChatNotFound) {
			WriteJSON(w, http.StatusNotFound, dto.NewErrorResponse(err))
			slog.Error("Failed to find chat",
				slog.Int("chatID", chatID),
				slog.String("error", err.Error()),
			)
			return
		}

		WriteJSON(w, http.StatusInternalServerError, dto.NewErrorResponse(err))
		slog.Error("Failed to find messages",
			slog.String("op", op),
			slog.Int("ChatID", chatID),
			slog.String("error", err.Error()),
		)
		return
	}

	WriteJSON(w, http.StatusOK, dto.NewMessagesResponse(chatWithMessages))
	slog.Info("Messages find successfuly",
		slog.String("op", op),
		slog.Int("chatID", chatID),
	)
}
