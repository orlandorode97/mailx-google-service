package users

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/orlandorode97/mailx-google-service/pkg/models"
)

type Endpoints struct {
	GetUserByIdEndpoint endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		GetUserByIdEndpoint: MakeGetUserByIdEndpoint(s),
	}
}

func MakeGetUserByIdEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, _ := request.(getUserByIdRequest)
		user, err := s.GetUserByID(req.UserID)
		if err != nil {
			return getUserByIdResponse{Err: err}, nil
		}
		return getUserByIdResponse{
			User: user,
		}, nil
	}
}

type getUserByIdRequest struct {
	UserID string
}

type getUserByIdResponse struct {
	Err  error        `json:"error,omitempty"`
	User *models.User `json:"user"`
}

func (g getUserByIdResponse) error() error {
	return g.Err
}
