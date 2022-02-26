package auth

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
)

func TestMakeHandler(t *testing.T) {
	t.Run("success - make handler returns the router", func(t *testing.T) {
		logger := log.NewLogfmtLogger(os.Stdin)
		auth := MockAuthService{}
		handler := MakeHandler(auth, logger)
		assert.NotNil(t, handler, "router is defined.")
	})
}

func TestDecodeCallbackRequest(t *testing.T) {
	t.Run("success - decodeCallbackRequest returns the request.", func(t *testing.T) {
		request, err := decodeCallbackRequest(context.Background(), &http.Request{})
		assert.Nil(t, err)
		assert.NotNil(t, request)
	})
}

func TestEncodeLoginResponse(t *testing.T) {
	t.Run("success - encodeLoginResponse returns the url response.", func(t *testing.T) {
		resp := loginResponse{
			AuthUrl: `http://localhost:3000/oauth=www.oauthurl.com/state=123456`,
		}
		expectedUrl := `{"auth_url":"http://localhost:3000/oauth=www.oauthurl.com/state=123456"}`
		w := httptest.NewRecorder()
		err := encodeLoginResponse(context.Background(), w, resp)
		assert.Nil(t, err)
		assert.Equal(t, expectedUrl, strings.TrimSuffix(w.Body.String(), "\n"))
	})
}

func TestDecodeCallBackRequest(t *testing.T) {
	t.Run("success - decodeCallbackRequest returns the request with state and code", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/login/callback", strings.NewReader("state=1234&code=1234"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.ParseForm()
		request, err := decodeCallbackRequest(context.Background(), req)
		r, ok := request.(callbackRequest)
		assert.Equal(t, true, ok)
		assert.Nil(t, err)
		assert.Equal(t, "1234", r.State)
		assert.Equal(t, "1234", r.Code)
	})
}

func TestEncodeCallBackRequest(t *testing.T) {
	testscases := []struct {
		name                string
		response            interface{}
		httpStatus          int
		redirectUrlExpected string
	}{
		{
			name: "success - redirects to the success url",
			response: callbackResponse{
				JWT: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			},
			httpStatus:          http.StatusPermanentRedirect,
			redirectUrlExpected: "http://localhost:3000/success?mailx_google_auth=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
		},
		{
			name: "failure - error configuring gmail service",
			response: callbackResponse{
				Err: errors.New("cannot configure gmail service."),
			},
			httpStatus:          http.StatusPermanentRedirect,
			redirectUrlExpected: "http://localhost:3000/error?error_message=cannot configure gmail service.",
		},
	}

	for _, test := range testscases {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			_ = encodeCallbackResponse(context.Background(), w, test.response)
			assert.Equal(t, test.redirectUrlExpected, w.Header().Get("Location"))
			assert.Equal(t, test.httpStatus, w.Result().StatusCode)
		})
	}
}
