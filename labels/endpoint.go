package labels

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateLabelEndpoint  endpoint.Endpoint
	DeleteLabelEndpoint  endpoint.Endpoint
	GetLabelByIdEndpoint endpoint.Endpoint
	GetLabelsEndpoint    endpoint.Endpoint
	UpdateLabelEndpoint  endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateLabelEndpoint:  MakeCreateLabelEndpoint(s),
		DeleteLabelEndpoint:  MakeDeleteLabelEndpoint(s),
		GetLabelByIdEndpoint: MakeGetLabelByIdEndpoint(s),
		GetLabelsEndpoint:    MakeGetLabelsEndpoint(s),
		UpdateLabelEndpoint:  MakeUpdateLabelEndpoint(s),
	}
}

func MakeGetUserProfileEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, nil
	}
}

func MakeCreateLabelEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, nil
	}
}

func MakeDeleteLabelEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, nil
	}
}

func MakeGetLabelByIdEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		token := ctx.Value(jwt.JWTContextKey).(string)
		fmt.Println(token)
		return nil, nil
	}
}

func MakeGetLabelsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, nil
	}
}

func MakeUpdateLabelEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, nil
	}
}
