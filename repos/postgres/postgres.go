package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/orlandorode97/mailx-google-service/repos"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) repos.Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) CreateUser(ctx context.Context, user *repos.User) error {
	query, args, err := sq.
		Insert("users").
		Columns("google_id", "name", "given_name", "family_name", "picture", "locale").
		Values(user.ID, user.Name, user.GivenName, user.FamilyName, user.Picture, user.Locale).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	if _, err = r.db.DB.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func (r *repository) GetUserByID(ctx context.Context, ID string) (*repos.User, error) {
	var user repos.User
	query, args, err := sq.
		Select("google_id", "name", "given_name", "family_name", "picture", "locale").
		From("users").
		Where("googl_id = $1", ID).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.DB.QueryRowContext(ctx, query, args...)
	if err = row.Scan(&user.ID, &user.Name, &user.GivenName, &user.FamilyName, &user.Picture, &user.Locale); err != nil {
		return nil, err
	}

	return &user, nil
}
