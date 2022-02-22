package labels

import (
	"context"
	"encoding/json"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/orlandorode97/mailx-google-service/pkg/middlewares"
	"github.com/orlandorode97/mailx-google-service/pkg/models"
)

// MakeHandler mounts the labels endpoints
func MakeHandler(labelService Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()

	e := MakeEndpoints(labelService)
	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(models.ErrorEncoder),
	}
	r.Methods(http.MethodPost).
		Path("/labels/").
		Handler(kithttp.NewServer(
			e.CreateLabelEndpoint,
			decodeLabelsRequest,
			encodeLabelsResponse,
			options...,
		))
	r.Methods(http.MethodDelete).
		Path("/labels/").
		Handler(kithttp.NewServer(
			e.DeleteLabelEndpoint,
			decodeLabelsRequest,
			encodeLabelsResponse,
			options...,
		))
	r.Methods(http.MethodGet).
		Path("/labels/{id:[0-9a-zA-Z\\W]+|}").
		Handler(kithttp.NewServer(
			e.GetLabelByIdEndpoint,
			decodeLabelsRequest,
			encodeLabelsResponse,
			options...,
		))
	r.Methods(http.MethodGet).
		Path("/labels/").
		Handler(kithttp.NewServer(
			e.GetLabelsEndpoint,
			decodeLabelsRequest,
			encodeLabelsResponse,
			options...,
		))
	r.Methods(http.MethodPut).
		Path("/labels/").
		Handler(kithttp.NewServer(
			e.UpdateLabelEndpoint,
			decodeLabelsRequest,
			encodeLabelsResponse,
			options...,
		))
	return middlewares.Authentication(r)
}

type labelsRequest struct{}

func decodeLabelsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if err, ok := r.Context().Value(middlewares.InvalidAuthKey).(error); ok {
		return nil, err
	}

	return labelsRequest{}, nil
}

func encodeLabelsResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
