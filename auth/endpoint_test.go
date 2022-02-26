package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/orlandorode97/mailx-google-service/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
)

type MockAuthService struct {
	mock.Mock
}

func (m MockAuthService) GetOauthUrl(ctx context.Context) (string, error) {
	args := m.Called(ctx)
	return args.String(0), args.Error(1)
}

func (m MockAuthService) GenerateOauthToken(ctx context.Context, code string) (*oauth2.Token, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (m MockAuthService) ConfigGmailServiceUser(ctx context.Context, code string) (*models.User, error) {
	args := m.Called(ctx, code)
	return args.Get(0).(*models.User), args.Error(1)
}
func (m MockAuthService) CreateJWT(ctx context.Context, user *models.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}

func TestMakeEndpoints(t *testing.T) {
	t.Run("success - MakeEndpoints returns a not nil auth endpoints.", func(t *testing.T) {
		var mockService MockAuthService
		endpoints := MakeEndpoints(mockService)
		assert.NotNil(t, endpoints, "endpoints is not nil.")
	})
}

func TestMakeGetOauthUrlEndpoint(t *testing.T) {
	testscases := []struct {
		name            string
		url             string
		err             error
		message         string
		endpointMessage string
	}{
		{
			name:            "success - endpoint returns authorization url.",
			url:             "http://wwww.randomurl.com",
			err:             nil,
			message:         "response is defined.",
			endpointMessage: "GetOauthUrlEndpoint is not nil.",
		},
		{
			name:            "failure - endpoint returns authorization url error.",
			url:             "",
			err:             errors.New("error after creating authorization url."),
			message:         "response is not defined.",
			endpointMessage: "GetOauthUrlEndpoint is not nil.",
		},
	}
	for _, test := range testscases {
		t.Run(test.name, func(t *testing.T) {
			var mockService MockAuthService
			ctx := context.Background()
			mockService.On("GetOauthUrl", ctx).Return(test.url, test.err)
			endpoint := MakeGetOauthUrlEndpoint(mockService)
			assert.NotNil(t, endpoint, test.endpointMessage)
			response, err := endpoint(ctx, loginRequest{})
			if err != nil {
				assert.Nil(t, response, test.message)
				assert.NotNil(t, err)
			}
			resp, ok := response.(loginResponse)
			if ok {
				assert.NotNil(t, resp, test.message)
			}
		})
	}
}

func TestMakeGetOauthCallbackEndpoint(t *testing.T) {
	testcases := []struct {
		name            string
		ctx             context.Context
		user            *models.User
		errGmailConfig  error
		jwt             string
		jwtError        error
		request         callbackRequest
		errorMessage    string
		responseMessage string
	}{
		{
			name: "success - user and jwt returned by the config gmail and create jwt functions",
			ctx:  context.Background(),
			user: &models.User{
				ID:        "1234567",
				Name:      "testing",
				GivenName: "testing again",
			},
			errGmailConfig: nil,
			jwt:            "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			jwtError:       nil,
			request: callbackRequest{
				Code:  "123456abcd",
				State: "alksasdkl1",
			},
			errorMessage:    "error from GetOauthCallbackEndpoint is nil.",
			responseMessage: "response is not nil.",
		},
		{
			name:           "failure - gmail config function returns an error",
			ctx:            context.Background(),
			user:           nil,
			errGmailConfig: errors.New("gmail config error due external dependency."),
			request: callbackRequest{
				Code:  "123456abcd",
				State: "alksasdkl1",
			},
			errorMessage:    "error from GetOauthCallbackEndpoint is not nil.",
			responseMessage: "response is nil.",
		},
		{
			name: "failure - error creating jwt.",
			ctx:  context.Background(),
			user: &models.User{
				ID:        "1234567",
				Name:      "testing",
				GivenName: "testing again",
			},
			errGmailConfig: nil,
			request: callbackRequest{
				Code:  "123456abcd",
				State: "alksasdkl1",
			},
			jwtError:        errors.New("error trying to generate the JWT."),
			errorMessage:    "error from GetOauthCallbackEndpoint is not nil.",
			responseMessage: "response is nil.",
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			var mockService MockAuthService
			mockService.On("ConfigGmailServiceUser", test.ctx, test.request.Code).Return(test.user, test.errGmailConfig)
			mockService.On("CreateJWT", test.ctx, test.user).Return(test.jwt, test.jwtError)
			endpoint := MakeGetOauthCallbackEndpoint(mockService)
			assert.NotNil(t, endpoint, "the endpoint GetOauthCallbackEndpoint is not nil.")
			response, err := endpoint(test.ctx, test.request)
			if err != nil {
				assert.Nil(t, response, test, test.responseMessage)
				assert.NotNil(t, err, test.errorMessage)
			}

			resp, ok := response.(callbackResponse)
			if ok {
				assert.NotNil(t, resp, test, test.responseMessage)
			}
		})
	}
}
