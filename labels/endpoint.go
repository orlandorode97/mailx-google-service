package labels

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"google.golang.org/api/gmail/v1"
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
		return nil, nil
	}
}

func MakeGetLabelsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getLabelsRequest)
		labels, err := s.GetLabels(req.UserID)
		if err != nil {
			return getLabelsResponse{Err: err}, nil
		}

		return getLabelsResponse{
			Labels: labels,
		}, nil
	}
}

func MakeUpdateLabelEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return nil, nil
	}
}

type getLabelsRequest struct {
	UserID string
}

type getLabelsResponse struct {
	Labels []*gmail.Label `json:"labels"`
	Err    error          `json:"error,omitempty"`
}

func (g getLabelsResponse) error() error {
	return g.Err
}
