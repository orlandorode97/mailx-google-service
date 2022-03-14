package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/orlandorode97/mailx-google-service/pkg/models"
)

func MakeHandler(authSvc Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeEndpoints(authSvc)
	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(models.ErrorEncoder),
	}

	r.Methods(http.MethodGet).
		Path("/auth/logout/").
		Handler(kithttp.NewServer(
			e.LogoutEndpoint,
			decodeLogoutRequest,
			encodeLogoutResponse,
			options...,
		))
	r.Methods(http.MethodGet).
		Path("/auth/login/").
		Handler(kithttp.NewServer(
			e.GetOauthUrlEndpoint,
			decodeLoginRequest,
			encodeLoginResponse,
			options...,
		))
	r.Methods(http.MethodGet).
		Path("/auth/callback/").
		Handler(kithttp.NewServer(
			e.GetOauthCallbackEndpoint,
			decodeCallbackRequest,
			encodeCallbackResponse,
			options...,
		))
	return r
}

func decodeLogoutRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return logoutRequest{}, nil
}

func encodeLogoutResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	cookie := &http.Cookie{
		Name:     "mailx_google_auth",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}

	w.Header().Add("Set-Cookie", cookie.String())
	return nil
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
	var redirectUrl string
	if e, ok := response.(errorer); ok && e.error() != nil {
		redirectUrl = fmt.Sprintf("http://localhost:3000/error?error_message=%s", e.error().Error())
		http.Redirect(w, &http.Request{}, redirectUrl, http.StatusPermanentRedirect)
		return nil
	}

	resp, _ := response.(callbackResponse)
	cookie := &http.Cookie{
		Name:     "mailx_google_auth",
		Value:    resp.JWT,
		Path:     "/",
		HttpOnly: true,
	}

	w.Header().Add("Set-Cookie", cookie.String())

	redirectUrl = "http://localhost:3000/success?mailx_google_success=true"
	http.Redirect(w, &http.Request{}, redirectUrl, http.StatusPermanentRedirect)
	return nil
}

type errorer interface {
	error() error
}
