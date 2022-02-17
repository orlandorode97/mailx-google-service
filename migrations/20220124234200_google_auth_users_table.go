package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upGoogleAuthUsersTable, downGoogleAuthUsersTable)
}

func upGoogleAuthUsersTable(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS auth_users(
			id SERIAL PRIMARY KEY,
			google_id CHAR(50) NOT NULL,
			access_token TEXT NOT NULL,
			token_expiration DATE NOT NULL,
			refresh_token TEXT NOT NULL,
			token_type CHAR(10) NOT NULL,
			is_active BOOLEAN NOT NULL DEFAULT TRUE,
			created_at TIMESTAMP NOT NULL DEFAULT now(),
			updated_at TIMESTAMP NOT NULL DEFAULT now()
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

func downGoogleAuthUsersTable(tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE IF EXISTS auth_users;")
	if err != nil {
		return err
	}
	return nil
}
