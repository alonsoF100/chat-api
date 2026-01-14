package postgres

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateMessages, downCreateMessages)
}

func upCreateMessages(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE messages (
			id serial PRIMARY KEY,
			chat_id INTEGER NOT NULL,
			text VARCHAR(5000) NOT NULL,
			created_at TIMESTAMP DEFAULT now(),
			CONSTRAINT fk_chat FOREIGN KEY (chat_id)
			REFERENCES chats(id) ON DELETE CASCADE
		);

		CREATE INDEX idx_messages_chat_id_created_at 
		ON messages(chat_id, created_at DESC);
	`)
	return err
}

func downCreateMessages(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE messages;")
	return err
}
