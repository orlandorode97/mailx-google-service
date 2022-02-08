package auth

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/orlandorode97/mailx-google-service/repos"
)

type Endpoints struct {
	GetOauthUrlEndpoint      endpoint.Endpoint
	GetOauthCallbackEndpoint endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		GetOauthUrlEndpoint:      MakeGetOauthUrlEndpoint(s),
		GetOauthCallbackEndpoint: MakeGetOauthCallbackEndpoint(s),
	}
}

func MakeGetOauthUrlEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		_ = request.(loginRequest)
		url, err := s.GetOauthUrl(ctx)
		if err != nil {
			return nil, err
		}
		return loginResponse{AuthUrl: url}, nil
	}
}
func MakeGetOauthCallbackEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(callbackRequest)
		token, err := s.GenerateOauthToken(ctx, req.Code)
		if err != nil {
			return callbackResponse{Err: err}, nil
		}

		client := s.GenerateHttpClient(ctx, token)
		err = s.ConfigGmailService(ctx, client)
		if err != nil {
			return callbackResponse{Err: err}, nil
		}

		err = s.CreateUser(ctx, token)
		if err != nil {
			return callbackResponse{Err: err}, nil
		}

		return callbackResponse{}, nil
	}
}

type loginRequest struct{}
type loginResponse struct {
	AuthUrl string `json:"auth_url"`
}

type callbackRequest struct {
	State string
	Code  string
}

type callbackResponse struct {
	User *repos.User `json:"user,omitempty"`
	Err  error       `json:"error,omitempty"`
}

func (c callbackResponse) error() error {
	return c.Err
}
