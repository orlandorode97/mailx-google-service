package auth

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/google/uuid"
	"github.com/orlandorode97/mailx-google-service/pkg/google"
	"github.com/orlandorode97/mailx-google-service/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

type MockMailxService struct {
	mock.Mock
}

func (m MockMailxService) GetGmailService(ID string) *gmail.Service {
	args := m.Called(ID)
	return args.Get(0).(*gmail.Service)
}

func (m MockMailxService) CreateGmailService(token *oauth2.Token) (*gmail.Service, error) {
	args := m.Called(token)
	return args.Get(0).(*gmail.Service), args.Error(1)
}

func (m MockMailxService) AddGmailServiceByID(ID string, gmailSvc *gmail.Service) {
	m.Called(ID, gmailSvc)
}

func (m MockMailxService) RecreateGmailService(ctx context.Context, ID string) (*gmail.Service, error) {
	args := m.Called(ctx, ID)
	return args.Get(0).(*gmail.Service), args.Error(1)
}

type MockOAuthConfig struct {
	mock.Mock
}

func (m MockOAuthConfig) Exchange(ctx context.Context, code string, opts ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	args := m.Called(ctx, code, opts)
	return args.Get(0).(*oauth2.Token), args.Error(1)
}

func (m MockOAuthConfig) AuthCodeURL(state string, opts ...oauth2.AuthCodeOption) string {
	args := m.Called(state, opts)
	return args.String(0)
}

type AuthServiceMock struct {
	config      google.OAuthConfiguration
	GetOauthUrl func(context.Context) (string, error)
	CreateJWT   func(context.Context, *models.User) (string, error)
}

type mockHttpTransport func(req *http.Request) *http.Response

func (m mockHttpTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m(req), nil
}

func NewTestClient(fn mockHttpTransport) *http.Client {
	return &http.Client{
		Transport: mockHttpTransport(fn),
	}
}

type MockDB struct {
	mock.Mock
}

func (db MockDB) CreateUser(ctx context.Context, user *models.User) error {
	args := db.Called(ctx, user)
	return args.Error(0)
}

func (db MockDB) GetUserByID(ctx context.Context, ID string) (*models.User, error) {
	args := db.Called(ctx, ID)
	return args.Get(0).(*models.User), args.Error(1)
}

func (db MockDB) GetTokenByUserId(ctx context.Context, ID string) (*models.Token, error) {
	args := db.Called(ctx, ID)
	return args.Get(0).(*models.Token), args.Error(1)
}

func (db MockDB) SaveAccessToken(ctx context.Context, ID string, token *oauth2.Token) error {
	args := db.Called(ctx, ID, token)
	return args.Error(0)
}

func (db MockDB) UpdateAccessToken(ctx context.Context, ID string, token *oauth2.Token) error {
	args := db.Called(ctx, ID, token)
	return args.Error(0)
}

func TestGetOauthUrl(t *testing.T) {
	testscases := []struct {
		name        string
		state       string
		expectedUrl string
		expectedErr error
		assertErr   func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
		assertUrl   func(t assert.TestingT, expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool
	}{
		{
			name:        "success - url is returned by TestGetOauthUrl",
			state:       uuid.NewString(),
			expectedUrl: "http://localhost:3000",
			expectedErr: nil,
			assertErr:   assert.Nil,
			assertUrl:   assert.Equal,
		},
		{
			name:        "failure - error returned by generating url by TestGetOauthUrl",
			state:       uuid.NewString(),
			expectedUrl: "",
			expectedErr: models.ErrAuthUrl{},
			assertErr:   assert.NotNil,
			assertUrl:   assert.Equal,
		},
	}

	for _, test := range testscases {
		t.Run(test.name, func(t *testing.T) {
			config := MockOAuthConfig{}
			config.On("AuthCodeURL", test.state, []oauth2.AuthCodeOption{oauth2.AccessTypeOffline}).Return(test.expectedUrl)
			svc := &AuthServiceMock{
				config: config,
			}
			svc.GetOauthUrl = func(c context.Context) (string, error) {
				url := svc.config.AuthCodeURL(test.state, oauth2.AccessTypeOffline)
				if url == "" {
					return "", models.ErrAuthUrl{}
				}
				return url, nil
			}
			url, err := svc.GetOauthUrl(context.Background())
			test.assertErr(t, err, test.expectedErr)
			test.assertUrl(t, test.expectedUrl, url)
		})
	}
}

func TestGenerateOauthToken(t *testing.T) {
	testcases := []struct {
		name          string
		code          string
		ctx           context.Context
		token         *oauth2.Token
		tokenErr      error
		assertErr     func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
		assertToken   func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
		expectedToken string
	}{
		{
			name: "success - oauth token generated.",
			code: "4/0AX4XfWhJiiRZlcRiBt3Q5a4Z5lg8A1P9VSc-edHoeRWDWtmyiYWKooMUh5JE6IKzzZXymw",
			ctx:  context.Background(),
			token: &oauth2.Token{
				AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			},
			tokenErr:    nil,
			assertErr:   assert.Nil,
			assertToken: assert.NotNil,
		},
		{
			name:        "failure - oauth token generation returns an error.",
			code:        "4/0AX4XfWhJiiRZlcRiBt3Q5a4Z5lg8A1P9VSc-edHoeRWDWtmyiYWKooMUh5JE6IKzzZXymw",
			ctx:         context.Background(),
			token:       nil,
			tokenErr:    errors.New("could not create oauth token access."),
			assertErr:   assert.NotNil,
			assertToken: assert.Nil,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			config := MockOAuthConfig{}
			config.On("Exchange", test.ctx, test.code, []oauth2.AuthCodeOption(nil)).Return(test.token, test.tokenErr)
			svc := &service{
				config: config,
				logger: log.NewLogfmtLogger(os.Stdin),
			}

			token, err := svc.GenerateOauthToken(test.ctx, test.code)
			test.assertErr(t, err)
			test.assertToken(t, token)

		})
	}
}

func TestConfigGmailServiceUser(t *testing.T) {
	testscases := []struct {
		name            string
		ctx             context.Context
		code            string
		token           *oauth2.Token
		tokenErr        error
		gmailSvc        *gmail.Service
		gmailSvcErr     error
		expectedUser    *models.User
		expectedUserErr error
		saveTokenErr    error
		assertErr       func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
		assertEqual     func(t assert.TestingT, expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool
	}{
		{
			name: "success - returns google user after configure a gmail service.",
			ctx:  context.Background(),
			code: "12345678",
			token: &oauth2.Token{
				AccessToken: "random access token",
			},
			gmailSvc: &gmail.Service{},
			expectedUser: &models.User{
				ID:   "1",
				Name: "Orlando",
			},
			assertErr:   assert.Nil,
			assertEqual: assert.Equal,
		},
		{
			name:         "failure - returns error when generating the oauth token access.",
			ctx:          context.Background(),
			code:         "12345678",
			token:        nil,
			tokenErr:     errors.New("error generating the oauth token access."),
			gmailSvc:     &gmail.Service{},
			expectedUser: &models.User{},
			assertErr:    assert.NotNil,
			assertEqual:  assert.NotEqual,
		},
		{
			name: "failure - returns error when configure gmail service",
			ctx:  context.Background(),
			code: "12345678",
			token: &oauth2.Token{
				AccessToken: "random access token",
			},
			gmailSvc:     &gmail.Service{},
			gmailSvcErr:  errors.New("error gmail service cannot be configured."),
			expectedUser: &models.User{},
			assertErr:    assert.NotNil,
			assertEqual:  assert.NotEqual,
		},
		{
			name: "failure - cannot create gmail user",
			ctx:  context.Background(),
			code: "12345678",
			token: &oauth2.Token{
				AccessToken: "random access token",
			},
			gmailSvc: &gmail.Service{},
			expectedUser: &models.User{
				ID:   "1",
				Name: "Orlando",
			},
			expectedUserErr: errors.New("error database is not running."),
			assertErr:       assert.NotNil,
			assertEqual:     assert.NotEqual,
		},
		{
			name: "failure - cannot save gmail access token",
			ctx:  context.Background(),
			code: "12345678",
			token: &oauth2.Token{
				AccessToken: "random access token",
			},
			gmailSvc: &gmail.Service{},
			expectedUser: &models.User{
				ID:   "1",
				Name: "Orlando",
			},
			saveTokenErr: errors.New("the table auth_users does not exists."),
			assertErr:    assert.NotNil,
			assertEqual:  assert.NotEqual,
		},
	}

	for _, test := range testscases {
		t.Run(test.name, func(t *testing.T) {
			config := MockOAuthConfig{}
			config.On("Exchange", test.ctx, test.code, []oauth2.AuthCodeOption(nil)).Return(test.token, test.tokenErr)

			db := MockDB{}
			db.On("GetUserByID", test.ctx, test.expectedUser.ID).Return(&models.User{}, test.expectedUserErr)
			db.On("CreateUser", test.ctx, test.expectedUser).Return(test.expectedUserErr)

			db.On("GetTokenByUserId", test.ctx, test.expectedUser.ID).Return(&models.Token{}, test.saveTokenErr)
			db.On("SaveAccessToken", test.ctx, test.expectedUser.ID, test.token).Return(test.saveTokenErr)
			db.On("UpdateAccessToken", test.ctx, test.expectedUser.ID, test.token).Return(test.saveTokenErr)

			mockMailxService := MockMailxService{}
			mockMailxService.On("CreateGmailService", test.token).Return(test.gmailSvc, test.gmailSvcErr)
			mockMailxService.On("AddGmailServiceByID", test.expectedUser.ID, test.gmailSvc)

			client := NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{"id": "1", "name": "Orlando"}`))),
				}
			})

			svc := &service{
				logger:       log.NewLogfmtLogger(os.Stdin),
				config:       config,
				db:           db,
				client:       client,
				mailxService: mockMailxService,
			}
			user, err := svc.ConfigGmailServiceUser(test.ctx, test.code)
			test.assertErr(t, err)
			test.assertEqual(t, test.expectedUser, user)
		})
	}
}

func TestCreateJWT(t *testing.T) {
	testcases := []struct {
		name        string
		ctx         context.Context
		token       string
		tokenErr    error
		user        *models.User
		assertErr   func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
		assertEqual func(t assert.TestingT, expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool
	}{
		{
			name:  "success - mailx token returned.",
			ctx:   context.Background(),
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			user: &models.User{
				ID: "1",
			},
			assertErr:   assert.Nil,
			assertEqual: assert.Equal,
		},
		{
			name:     "failure - cannot find JWT_SIGNING_KEY",
			ctx:      context.Background(),
			token:    "",
			tokenErr: errors.New("JWT_SIGNING_KEY is missing"),
			user: &models.User{
				ID: "1",
			},
			assertErr:   assert.NotNil,
			assertEqual: assert.Equal,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			svc := AuthServiceMock{}
			svc.CreateJWT = func(c context.Context, u *models.User) (string, error) {
				return test.token, test.tokenErr
			}

			token, err := svc.CreateJWT(test.ctx, test.user)
			test.assertErr(t, err)
			test.assertEqual(t, test.token, token)
		})
	}
}

func TestCreateUser(t *testing.T) {
	testscases := []struct {
		name           string
		ctx            context.Context
		body           string
		expectedUser   *models.User
		getUserByIdErr error
		createUserErr  error
		assertErr      func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
		assertUser     func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
	}{
		{
			name:           "success - user created in database",
			ctx:            context.Background(),
			body:           `{"id": "12345", "name": "Orlando"}`,
			getUserByIdErr: sql.ErrNoRows,
			expectedUser: &models.User{
				ID:   "12345",
				Name: "Orlando",
			},
			assertErr:  assert.Nil,
			assertUser: assert.NotNil,
		},
		{
			name: "success - user already exists, skipping creation ",
			ctx:  context.Background(),
			body: `{"id": "12345", "name": "Orlando"}`,
			expectedUser: &models.User{
				ID:   "12345",
				Name: "Orlando",
			},
			assertErr:  assert.Nil,
			assertUser: assert.NotNil,
		},
		{
			name:           "failure - cannot get user by id before its creation.",
			ctx:            context.Background(),
			body:           `{"id": "12345", "name": "Orlando"}`,
			getUserByIdErr: errors.New("the database is not connected."),
			expectedUser: &models.User{
				ID:   "12345",
				Name: "Orlando",
			},
			assertErr:  assert.NotNil,
			assertUser: assert.Nil,
		},
		{
			name:           "failure - cannot create user",
			ctx:            context.Background(),
			body:           `{"id": "12345", "name": "Orlando"}`,
			getUserByIdErr: sql.ErrNoRows,
			createUserErr:  errors.New("cannot insert a boolean value in user table"),
			expectedUser: &models.User{
				ID:   "12345",
				Name: "Orlando",
			},
			assertErr:  assert.NotNil,
			assertUser: assert.Nil,
		},
	}

	for _, test := range testscases {
		t.Run(test.name, func(t *testing.T) {
			client := NewTestClient(func(req *http.Request) *http.Response {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(test.body))),
				}
			})
			db := MockDB{}
			db.On("GetUserByID", test.ctx, test.expectedUser.ID).Return(&models.User{}, test.getUserByIdErr)
			db.On("CreateUser", test.ctx, test.expectedUser).Return(test.createUserErr)
			logger := log.NewLogfmtLogger(os.Stdin)
			svc := service{
				db:     db,
				logger: logger,
				client: client,
			}
			user, err := svc.createUser(test.ctx, &oauth2.Token{})
			test.assertErr(t, err)
			test.assertUser(t, user)
		})
	}
}

func TestSaveAccessToken(t *testing.T) {
	testcases := []struct {
		name                 string
		ctx                  context.Context
		user                 *models.User
		token                *oauth2.Token
		mailxToken           *models.Token
		getTokenByUserIdErr  error
		saveAccessTokenErr   error
		updateAccessTokenErr error
		assertErr            func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
	}{
		{
			name: "success - oauth access token does not exist in db, create a new entry",
			ctx:  context.Background(),
			user: &models.User{
				ID: "1",
			},
			token: &oauth2.Token{
				AccessToken:  "token access",
				RefreshToken: "refresh token access",
				TokenType:    "Bearer",
			},
			mailxToken: &models.Token{
				AccessToken:  "token access",
				RefreshToken: "refresh token access",
				TokenType:    "Bearer",
			},
			getTokenByUserIdErr: sql.ErrNoRows,
			assertErr:           assert.Nil,
		},
		{
			name: "failure - cannot create access token entry in db",
			ctx:  context.Background(),
			user: &models.User{
				ID: "1",
			},
			token: &oauth2.Token{
				AccessToken:  "token access",
				RefreshToken: "refresh token access",
				TokenType:    "Bearer",
			},
			mailxToken: &models.Token{
				AccessToken:  "token access",
				RefreshToken: "refresh token access",
				TokenType:    "Bearer",
			},
			getTokenByUserIdErr: sql.ErrNoRows,
			saveAccessTokenErr:  errors.New("database is not running"),
			assertErr:           assert.NotNil,
		},
		{
			name: "failure - cannot update access token entry in db",
			ctx:  context.Background(),
			user: &models.User{
				ID: "1",
			},
			token: &oauth2.Token{
				AccessToken:  "token access",
				RefreshToken: "refresh token access",
				TokenType:    "Bearer",
			},
			mailxToken: &models.Token{
				AccessToken:  "token access",
				RefreshToken: "refresh token access",
				TokenType:    "Bearer",
			},
			assertErr:            assert.NotNil,
			updateAccessTokenErr: errors.New("cannot update oauth token access"),
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			db := MockDB{}
			db.On("GetTokenByUserId", test.ctx, test.user.ID).Return(test.mailxToken, test.getTokenByUserIdErr)
			db.On("SaveAccessToken", test.ctx, test.user.ID, test.token).Return(test.saveAccessTokenErr)
			db.On("UpdateAccessToken", test.ctx, test.user.ID, test.token).Return(test.updateAccessTokenErr)
			logger := log.NewLogfmtLogger(os.Stdin)
			svc := service{
				db:     db,
				logger: logger,
			}
			err := svc.saveAccessToken(test.ctx, test.user.ID, test.token)
			test.assertErr(t, err)
		})
	}
}
