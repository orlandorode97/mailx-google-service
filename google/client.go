package google

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
)

func NewClient(config *oauth2.Config) (*http.Client, error) {
	ctx := context.Background()

	url := config.AuthCodeURL("random-thing", oauth2.AccessTypeOffline)
	fmt.Println("url", url)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, err
	}
	token, err := config.Exchange(ctx, authCode)
	if err != nil {
		return nil, err
	}

	return config.Client(context.Background(), token), nil
}
