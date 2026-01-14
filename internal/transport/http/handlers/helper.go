package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	const op = "handlers/WriteJSON"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Debug("Failed to encode",
			slog.String("op", op),
			slog.Any("data", data),
			slog.String("error", err.Error()),
		)
	}
}
