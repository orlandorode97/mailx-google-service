package messages

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/orlandorode97/mailx-google-service/pkg/models"
)

type Endpoints struct {
	GetMessagesEndpoint    endpoint.Endpoint
	GetMessageByIDEndpoint endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		GetMessagesEndpoint:    MakeGetMessages(s),
		GetMessageByIDEndpoint: MakeGetMessageByID(s),
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

func MakeGetMessageByID(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getMessageByIDRequest)
		message, err := s.GetMessageByID(ctx, req.UserID, req.MessageID)
		if err != nil {
			return getMessageByIDResponse{
				Err: err,
			}, nil
		}
		return getMessageByIDResponse{
			Message: message,
		}, nil
	}
}

type getMessagesRequest struct {
	UserID string
}

type getMessagesResponse struct {
	Messages []*models.Message `json:"messages"`
	Err      error             `json:"error,omitempty"`
}

type getMessageByIDRequest struct {
	UserID    string
	MessageID string
}

type getMessageByIDResponse struct {
	Message *models.Message `json:"message"`
	Err     error           `json:"error,omitempty"`
}
