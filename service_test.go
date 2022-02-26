package mailx

import (
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

type MockMailxService struct {
	mock.Mock
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
