package models

import (
	"context"
	"encoding/json"
	"net/http"
)

type ErrAuthUrl struct{}

func (e ErrAuthUrl) Error() string {
	return "there is problem trying to generate the oauth url."
}

type ErrInvalidCookie struct{}

func (e ErrInvalidCookie) Error() string {
	return "the cookie is not presented or is invalid."
}

type ErrInvalidToken struct{}

func (e ErrInvalidToken) Error() string {
	return "the token is invalid."
}

type ErrExpiredToken struct{}

func (e ErrExpiredToken) Error() string {
	return "the token has been expired."
}

type ErrMalformedToken struct{}

func (e ErrMalformedToken) Error() string {
	return "the token is malformed."
}

type ErrInactiveToken struct{}

func (e ErrInactiveToken) Error() string {
	return "the token is inactive."
}

// ErrorEncoder encodes incoming errors to write the corresponding http status header.
func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {

	switch err.(type) {
	case ErrAuthUrl:
		w.WriteHeader(http.StatusServiceUnavailable)
	case ErrInvalidToken, ErrExpiredToken, ErrMalformedToken, ErrInactiveToken, ErrInvalidCookie:
		w.WriteHeader(http.StatusUnauthorized)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
