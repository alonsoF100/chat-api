package apperrors

import "errors"

// TODO Добавить константные ошибки

var (
	ErrChatNotFound = errors.New("chat not found")
)
