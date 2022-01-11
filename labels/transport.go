package labels

import (
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

// MakeHandler mounts the labels endpoints
func MakeHandler(labelService Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeEndpoints(labelService)
	options := []kithttp.ServerOption{}
	r.Methods(http.MethodPost).
		Path("/labels/").
		Handler(kithttp.NewServer(
			e.CreateLabelEndpoint,
			nil,
			nil,
			options...,
		))
	r.Methods(http.MethodDelete).
		Path("/labels/").
		Handler(kithttp.NewServer(
			e.DeleteLabelEndpoint,
			nil,
			nil,
			options...,
		))
	r.Methods(http.MethodGet).
		Path("/labels/{id:[0-9a-zA-Z\\W]+|}").
		Handler(kithttp.NewServer(
			e.GetLabelByIdEndpoint,
			nil,
			nil,
			options...,
		))
	r.Methods(http.MethodGet).
		Path("/labels/").
		Handler(kithttp.NewServer(
			e.GetLabelsEndpoint,
			nil,
			nil,
			options...,
		))
	r.Methods(http.MethodPut).
		Path("/labels/").
		Handler(kithttp.NewServer(
			e.UpdateLabelEndpoint,
			nil,
			nil,
			options...,
		))
	return r
}
