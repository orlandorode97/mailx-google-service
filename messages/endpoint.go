package messages

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"google.golang.org/api/gmail/v1"
)

type Endpoints struct {
	GetMessagesEndpoint endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		GetMessagesEndpoint: MakeGetMessages(s),
	}
}

func MakeGetMessages(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getMessagesRequest)
		messages, err := s.GetMessages(ctx, req.UserID)
		if err != nil {
			return getMessagesResponse{
				Err: err,
			}, nil
		}
		return getMessagesResponse{
			Messages: messages,
		}, nil
	}
}

type getMessagesRequest struct {
	UserID string
}

type getMessagesResponse struct {
	Messages []*gmail.Message `json:"messages"`
	Err      error            `json:"error,omitempty"`
}
