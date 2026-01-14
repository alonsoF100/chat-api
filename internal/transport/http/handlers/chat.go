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
pattern: /chats/
method: POST
info: JSON in request body

succeed:

	-status code: 201 created
	-response body: JSON represented created chat

failed:

	-status code: 400 bad request, 500 internal server error
	-response body: JSON with error message + timestamp
*/
func (h Handler) CreateChat(w http.ResponseWriter, r *http.Request) {
	const op = "handlers/CreateChat"

	var req dto.CreateChatRequest
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

	chat, err := h.ChatService.CreateChat(req.Title)
	if err != nil {
		WriteJSON(w, http.StatusInternalServerError, dto.NewErrorResponse(err))
		slog.Debug("Failed to create chat",
			slog.String("op", op),
			slog.String("error", err.Error()),
		)
		return
	}

	WriteJSON(w, http.StatusCreated, dto.NewChatResponse(chat))
	slog.Info("Chat created successfuly",
		slog.Int("chatID", chat.ID),
	)
}

/*
pattern: /chats/{id}
method: DELETE
info: chatID from URL

succeed:

	-status code: 204 no content
	-response body: JSON represented chat + []messages

failed:

	-status code: 400 bad request, 404 not found, 500 internal server error
	-response body: JSON with error message + timestamp
*/
func (h Handler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	const op = "handlers/DeleteChat"

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

	err = h.ChatService.DeleteChat(chatID)
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
		slog.Error("Failed to delete chat",
			slog.String("op", op),
			slog.Int("ChatID", chatID),
			slog.String("error", err.Error()),
		)
		return
	}

	WriteJSON(w, http.StatusNoContent, nil)
	slog.Info("Chat deleted successfuly",
		slog.String("op", op),
		slog.Int("chatID", chatID),
	)
}
