package users

import (
	"context"
	"encoding/json"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/orlandorode97/mailx-google-service/pkg/middlewares"
	"github.com/orlandorode97/mailx-google-service/pkg/models"
)

func MakeHandler(usersService Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeEndpoints(usersService)

	options := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(models.ErrorEncoder),
	}

	r.Methods(http.MethodGet).
		Path("/users/").
		Handler(kithttp.NewServer(
			e.GetUserByIdEndpoint,
			decodeUsersRequest,
			encodeUsersResponse,
			options...,
		))

	return middlewares.Authentication(r)
}

func decodeUsersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	if err, ok := r.Context().Value(middlewares.InvalidAuthKey).(error); ok && err != nil {
		return nil, err
	}

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		return nil, models.ErrInvalidData{Field: "user_id"}
	}

	return getUserByIdRequest{
		UserID: userID,
	}, nil
}

func encodeUsersResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		return e.error()
	}

	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}
