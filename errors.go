package mailx

import (
	"context"
	"encoding/json"
	"net/http"
)

type ErrAuthUrl struct{}

func (e ErrAuthUrl) Error() string {
	return "there is problem trying to generate the oauth url."
}

// ErrorEncoder encodes incoming errors to write the corresponding http status header.
func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	switch err.(type) {
	case ErrAuthUrl:
		w.WriteHeader(http.StatusServiceUnavailable)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
