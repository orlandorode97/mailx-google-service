package google

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

func NewClient(config *oauth2.Config, token *oauth2.Token) *http.Client {
	return config.Client(context.Background(), token)
}
