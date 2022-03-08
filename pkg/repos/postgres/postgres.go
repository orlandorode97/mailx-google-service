package postgres

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/orlandorode97/mailx-google-service/pkg/models"
	"github.com/orlandorode97/mailx-google-service/pkg/repos"
	"golang.org/x/oauth2"
)

type repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) repos.Repository {
	if err := db.Ping(); err != nil {
		return nil
	}

	return &repository{
		db: db,
	}
}

func (r *repository) CreateUser(ctx context.Context, user *models.User) error {
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

func (r *repository) GetUserByID(ctx context.Context, ID string) (*models.User, error) {
	var user models.User
	query, args, err := sq.
		Select("google_id", "name", "given_name", "family_name", "picture", "locale").
		From("users").
		Where(sq.Eq{"google_id": ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.DB.QueryRowContext(ctx, query, args...)
	if err = row.Scan(
		&user.ID,
		&user.Name,
		&user.GivenName,
		&user.FamilyName,
		&user.Picture,
		&user.Locale,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetTokenByUserId(ctx context.Context, ID string) (*models.Token, error) {
	var token models.Token
	query, args, err := sq.
		Select("id", "google_id", "access_token", "token_expiration", "refresh_token", "token_type").
		From("auth_users").
		Where(sq.Eq{"google_id": ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}
	row := r.db.DB.QueryRowContext(ctx, query, args...)
	if err := row.Scan(
		&token.ID,
		&token.UserID,
		&token.AccessToken,
		&token.TokenExpiration,
		&token.RefreshToken,
		&token.TokenType,
	); err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *repository) SaveAccessToken(ctx context.Context, ID string, token *oauth2.Token) error {
	query, args, err := sq.
		Insert("auth_users").
		Columns("google_id", "access_token", "token_expiration", "refresh_token", "token_type").
		Values(ID, token.AccessToken, token.Expiry, token.RefreshToken, token.TokenType).
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

func (r *repository) UpdateAccessToken(ctx context.Context, ID string, token *oauth2.Token) error {
	query, args, err := sq.
		Update("auth_users").
		Set("access_token", token.AccessToken).
		Set("token_expiration", token.Expiry).
		Set("refresh_token", token.RefreshToken).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"google_id": ID}).
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
