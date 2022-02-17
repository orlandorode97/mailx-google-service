package auth

import (
	"context"

	"github.com/go-kit/kit/endpoint"
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
		user, err := s.ConfigGmailServiceUser(ctx, req.Code)
		if err != nil {
			return callbackResponse{Err: err}, nil
		}
		jwt, err := s.CreateJWT(ctx, user)
		if err != nil {
			return callbackResponse{Err: err}, nil
		}
		return callbackResponse{JWT: jwt}, nil
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
	JWT string `json:"-"`
	Err error  `json:"error,omitempty"`
}

func (c callbackResponse) error() error {
	return c.Err
}
