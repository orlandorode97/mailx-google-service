package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/orlandorode97/mailx-google-service"
)

func MakeHandler(authSvc Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeEndpoints(authSvc)
	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(mailx.ErrorEncoder),
	}
	r.Methods(http.MethodGet).
		Path("/login/").
		Handler(kithttp.NewServer(
			e.GetOauthUrlEndpoint,
			decodeLoginRequest,
			encodeLoginResponse,
			options...,
		))
	r.Methods(http.MethodGet).
		Path("/login/callback/").
		Handler(kithttp.NewServer(
			e.GetOauthCallbackEndpoint,
			decodeCallbackRequest,
			encodeCallbackResponse,
			options...,
		))
	return r
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return loginRequest{}, nil
}

func encodeLoginResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	// check for possible errors coming form the response
	return json.NewEncoder(w).Encode(response)
}

func decodeCallbackRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return callbackRequest{
		State: r.FormValue("state"),
		Code:  r.FormValue("code"),
	}, nil
}

func encodeCallbackResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		redirectUrl := fmt.Sprintf("http://localhost:3000/error?error_message=%s", e.error().Error())
		http.Redirect(w, &http.Request{}, redirectUrl, http.StatusPermanentRedirect)
		return nil
	}
	http.Redirect(w, &http.Request{}, "http://localhost:3000/inbox", http.StatusTemporaryRedirect)
	return nil
}

type errorer interface {
	error() error
}