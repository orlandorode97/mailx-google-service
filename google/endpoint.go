package google

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

type Endpoints struct {
	GetUserProfileEndpoint endpoint.Endpoint
	CreateLabelEndpoint    endpoint.Endpoint
	DeleteLabelEndpoint    endpoint.Endpoint
	GetLabelByIdEndpoint   endpoint.Endpoint
	GetLabelsEndpoint      endpoint.Endpoint
	UpdateLabelEndpoint    endpoint.Endpoint
}

func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		GetUserProfileEndpoint: MakeGetUserProfileEndpoint(s),
		CreateLabelEndpoint:    MakeCreateLabelEndpoint(s),
		DeleteLabelEndpoint:    MakeDeleteLabelEndpoint(s),
		GetLabelByIdEndpoint:   MakeGetLabelByIdEndpoint(s),
		GetLabelsEndpoint:      MakeGetLabelsEndpoint(s),
		UpdateLabelEndpoint:    MakeUpdateLabelEndpoint(s),
	}
}

func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = fmt.Sprintf("https://%s", instance)
	}
	target, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, nil
	}
	options := []httptransport.ClientOption{}

	return Endpoints{
		GetUserProfileEndpoint: httptransport.NewClient(http.MethodGet, target, nil, nil, options...).Endpoint(),
		CreateLabelEndpoint:    httptransport.NewClient(http.MethodPost, target, nil, nil, options...).Endpoint(),
		DeleteLabelEndpoint:    httptransport.NewClient(http.MethodDelete, target, nil, nil, options...).Endpoint(),
		GetLabelByIdEndpoint:   httptransport.NewClient(http.MethodGet, target, nil, nil, options...).Endpoint(),
		GetLabelsEndpoint:      httptransport.NewClient(http.MethodGet, target, nil, nil, options...).Endpoint(),
		UpdateLabelEndpoint:    httptransport.NewClient(http.MethodPut, target, nil, nil, options...).Endpoint(),
	}, nil
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
