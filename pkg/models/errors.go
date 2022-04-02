package models

import (
	"context"
	"encoding/json"
	"fmt"
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

type ErrInvalidSignature struct{}

func (e ErrInvalidSignature) Error() string {
	return "the token has an invalid signature."
}

type ErrInvalidData struct {
	Field string
}

func (e ErrInvalidData) Error() string {
	return fmt.Sprintf("the field `%s` is invalid", e.Field)
}

// ErrorEncoder encodes incoming errors to write the corresponding http status header.
func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {

	switch err.(type) {
	case ErrInvalidData:
		w.WriteHeader(http.StatusBadRequest)
	case ErrAuthUrl:
		w.WriteHeader(http.StatusServiceUnavailable)
	case ErrInvalidSignature, ErrInvalidToken, ErrExpiredToken, ErrMalformedToken, ErrInactiveToken, ErrInvalidCookie:
		w.WriteHeader(http.StatusUnauthorized)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
