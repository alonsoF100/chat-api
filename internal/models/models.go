package models

import "time"

type Chat struct {
	ID        int
	Title     string
	CreatedAt time.Time
}

type Message struct {
	ID        int
	ChatID    int
	Text      string
	CreatedAt time.Time
}

type ChatWithMessages struct {
	Chat     Chat
	Messages []*Message
}
