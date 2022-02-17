package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upGoogleUsersTableGo, downGoogleUsersTableGo)
}

func upGoogleUsersTableGo(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS users(
			google_id CHAR(50) PRIMARY KEY,
			name CHAR(50) NOT NULL,
			given_name CHAR(50) NOT NULL,
			family_name CHAR(50) NOT NULL,
			picture CHAR(100) NOT NULL,
			locale CHAR(10),
			created_at TIMESTAMP NOT NULL DEFAULT now(),
			updated_at TIMESTAMP NOT NULL DEFAULT now()
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

func downGoogleUsersTableGo(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`DROP TABLE IF EXISTS users;`)
	if err != nil {
		return err
	}
	return nil
}
