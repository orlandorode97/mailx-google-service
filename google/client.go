package google

import (
	"context"

	"google.golang.org/api/gmail/v1"
)

func New() (*gmail.Service, error) {
	ctx := context.Background()
	gmailService, err := gmail.NewService(ctx)
	if err != nil {
		return nil, err
	}
	return gmailService, nil
}
