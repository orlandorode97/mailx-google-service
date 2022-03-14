package mailx

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/orlandorode97/mailx-google-service/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

type MockMailxService struct {
	mock.Mock
}

type MockDB struct {
	mock.Mock
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

func TestAddGmailServiceByID(t *testing.T) {
	t.Run("success - gmail is added by user ID", func(t *testing.T) {
		logger := log.NewLogfmtLogger(os.Stdout)
		svc := &service{
			logger:    logger,
			gmailSvcs: make(map[string]*gmail.Service),
		}
		svc.AddGmailServiceByID("1", &gmail.Service{})
		assert.Len(t, svc.gmailSvcs, 1)
	})
}

func TestGetGmailService(t *testing.T) {
	testscases := []struct {
		name           string
		userId         string
		gmailSvcs      map[string]*gmail.Service
		assertGmailSvc func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
	}{
		{
			name:   "success - gmail service is returned by user id",
			userId: "1",
			gmailSvcs: map[string]*gmail.Service{
				"1": {
					BasePath:  "",
					UserAgent: "",
					Users:     &gmail.UsersService{},
				},
			},
			assertGmailSvc: assert.NotNil,
		},
		{
			name:   "failure - cannot find gmail service attached by user id ",
			userId: "1",
			gmailSvcs: map[string]*gmail.Service{
				"2": {
					BasePath:  "",
					UserAgent: "",
					Users:     &gmail.UsersService{},
				},
			},
			assertGmailSvc: assert.Nil,
		},
	}

	for _, test := range testscases {
		t.Run(test.name, func(t *testing.T) {
			logger := log.NewLogfmtLogger(os.Stdout)
			svc := &service{
				logger:    logger,
				gmailSvcs: test.gmailSvcs,
			}
			gmailSvc := svc.GetGmailService(test.userId)
			test.assertGmailSvc(t, gmailSvc)
		})
	}
}

func TestCreateGmailService(t *testing.T) {
	t.Run("success - gmail service is created", func(t *testing.T) {
		svc := &service{
			logger: log.NewLogfmtLogger(os.Stdout),
			config: &oauth2.Config{},
		}
		token := &oauth2.Token{}
		gmailSvc, err := svc.CreateGmailService(token)
		assert.Nil(t, err)
		assert.NotNil(t, gmailSvc)
	})
}

func TestRecreateGmailService(t *testing.T) {
	testcases := []struct {
		name      string
		ctx       context.Context
		userId    string
		token     *models.Token
		errToken  error
		errSvc    error
		assertErr func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
		assertSvc func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
	}{
		{
			name:   "success - gmail service recreated",
			ctx:    context.Background(),
			userId: "1",
			token: &models.Token{
				AccessToken:     "access-token",
				RefreshToken:    "refresh-token",
				TokenType:       "Bearer",
				TokenExpiration: time.Now(),
			},
			assertErr: assert.Nil,
			assertSvc: assert.NotNil,
		},
		{
			name:      "failure - get token by id fails",
			ctx:       context.Background(),
			userId:    "1",
			token:     nil,
			errToken:  errors.New("database is not online"),
			assertErr: assert.NotNil,
			assertSvc: assert.Nil,
		},
		{
			name:   "failure - gmail cannot be recreated",
			ctx:    context.Background(),
			userId: "1",
			token: &models.Token{
				AccessToken:     "access-token",
				RefreshToken:    "refresh-token",
				TokenType:       "Bearer",
				TokenExpiration: time.Now(),
			},
			errSvc:    errors.New("token no valid"),
			assertErr: assert.Nil,
			assertSvc: assert.NotNil,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			db := MockDB{}
			db.On("GetTokenByUserId", test.ctx, test.userId).Return(test.token, test.errToken)
			logger := log.NewLogfmtLogger(os.Stdout)
			svc := New(logger, db, nil)
			newSvc, err := svc.RecreateGmailService(test.ctx, test.userId)
			test.assertErr(t, err)
			test.assertSvc(t, newSvc)
		})
	}
}
