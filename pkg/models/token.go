package models

import "time"

type Token struct {
	ID              string
	UserID          string
	AccessToken     string
	TokenExpiration time.Time
	RefreshToken    string
	TokenType       string
}
