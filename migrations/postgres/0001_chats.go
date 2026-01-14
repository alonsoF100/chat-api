package postgres

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateChats, downCreateChats)
}

func upCreateChats(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		CREATE TABLE chats (
			id serial PRIMARY KEY,
			title VARCHAR(200) NOT NULL,
			created_at TIMESTAMP DEFAULT now()
		);
	`)
	return err
}

func downCreateChats(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE chats;")
	return err
}
