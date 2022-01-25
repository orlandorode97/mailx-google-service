package migrations

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigration(upAuthUsersRelation, downAuthUsersRelation)
}

func upAuthUsersRelation(tx *sql.Tx) error {
	_, err := tx.Exec(`
		ALTER TABLE auth_users
		ADD CONSTRAINT auth_users_fk FOREIGN KEY (google_id) REFERENCES users (google_id) MATCH FULL;
	`)
	if err != nil {
		return err
	}
	return nil
}

func downAuthUsersRelation(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	_, err := tx.Exec(`
		ALTER TABLE IF EXISTS auth_users
		DROP CONSTRAINT IF EXISTS auth_users_fk;`,
	)
	if err != nil {
		return err
	}
	return nil
}
