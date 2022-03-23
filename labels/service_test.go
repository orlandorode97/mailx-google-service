package labels

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/go-kit/log"
	"github.com/orlandorode97/mailx-google-service/pkg/google"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
)

type MockGmailService struct {
	mock.Mock
}

func (m MockGmailService) GetLabelsService() google.Labeler {
	args := m.Called()
	return args.Get(0).(google.Labeler)
}
func (m MockGmailService) GetMessagesService() google.Messenger {
	args := m.Called()
	return args.Get(0).(google.Messenger)
}

type MockLabeler struct {
	mock.Mock
}

func (m MockLabeler) Create(ID string, label *gmail.Label) google.LabelerClient {
	args := m.Called(ID, label)
	return args.Get(0).(google.LabelerClient)
}
func (m MockLabeler) Delete(userID string, labelID string) google.LabelerClientDelete {
	args := m.Called(userID, labelID)
	return args.Get(0).(google.LabelerClientDelete)
}
func (m MockLabeler) Get(userID string, labelID string) google.LabelerClient {
	args := m.Called(userID, labelID)
	return args.Get(0).(google.LabelerClient)
}
func (m MockLabeler) List(userID string) google.LabelerClientList {
	args := m.Called(userID)
	return args.Get(0).(google.LabelerClientList)
}
func (m MockLabeler) Patch(userID string, labelID string, label *gmail.Label) google.LabelerClient {
	args := m.Called(userID, labelID, label)
	return args.Get(0).(google.LabelerClient)
}
func (m MockLabeler) Update(userID string, labelID string, label *gmail.Label) google.LabelerClient {
	args := m.Called(userID, labelID, label)
	return args.Get(0).(google.LabelerClient)
}

type MockMailxService struct {
	mock.Mock
}

func (m MockMailxService) GetGmailService(userID string) google.Service {
	args := m.Called(userID)
	return args.Get(0).(google.Service)
}

func (m MockMailxService) CreateGmailService(token *oauth2.Token) (google.Service, error) {
	args := m.Called(token)
	return args.Get(0).(google.Service), args.Error(1)
}

func (m MockMailxService) AddGmailServiceByID(ID string, gmailSvc google.Service) google.Service {
	args := m.Called(ID, gmailSvc)
	return args.Get(0).(google.Service)
}

func (m MockMailxService) RecreateGmailService(ctx context.Context, ID string) (google.Service, error) {
	args := m.Called(ctx, ID)
	return args.Get(0).(google.Service), args.Error(1)
}

type MockLabelerClientList struct {
	mock.Mock
}

func (m MockLabelerClientList) Do(opts ...googleapi.CallOption) (*gmail.ListLabelsResponse, error) {
	args := m.Called(opts)
	return args.Get(0).(*gmail.ListLabelsResponse), args.Error(1)
}

func TestGetLabels(t *testing.T) {
	testcases := []struct {
		name                string
		ctx                 context.Context
		userID              string
		gmailSvcRecreate    google.Service
		isGmailSvcNil       bool
		errGmailSvcRecreate error
		labelResponse       *gmail.ListLabelsResponse
		errLabels           error
		assertErr           func(t assert.TestingT, object interface{}, msgAndArgs ...interface{}) bool
		assertEqual         func(t assert.TestingT, expected interface{}, actual interface{}, msgAndArgs ...interface{}) bool
	}{
		{
			name:   "success - labels returned by the gmail api.",
			ctx:    context.Background(),
			userID: "1",
			labelResponse: &gmail.ListLabelsResponse{
				Labels: []*gmail.Label{
					{
						Name: "Label 1",
						Id:   "LABEL_1",
					},

					{
						Name: "Label 2",
						Id:   "LABLE_2",
					},
				},
			},
			assertErr:   assert.Nil,
			assertEqual: assert.Equal,
		},
		{
			name:                "failure - cannot recreate gmail service and cannot return gmail labels service.",
			ctx:                 context.Background(),
			userID:              "1",
			isGmailSvcNil:       true,
			errGmailSvcRecreate: errors.New("Cannot recreate gmail service"),
			labelResponse:       &gmail.ListLabelsResponse{},
			assertErr:           assert.NotNil,
			assertEqual:         assert.Equal,
		},
		{
			name:          "failure - gmail labels service responds an error.",
			ctx:           context.Background(),
			userID:        "1",
			labelResponse: &gmail.ListLabelsResponse{},
			errLabels:     errors.New("request timeout"),
			assertErr:     assert.NotNil,
			assertEqual:   assert.Equal,
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			t.Run(test.name, func(t *testing.T) {
				logger := log.NewLogfmtLogger(os.Stdin)
				mockGmailService := MockGmailService{}
				mailxSvc := MockMailxService{}
				mockLabeler := MockLabeler{}
				mockCall := MockLabelerClientList{}

				mockCall.On("Do", []googleapi.CallOption(nil)).Return(test.labelResponse, test.errLabels)
				mockLabeler.On("List", test.userID).Return(mockCall)
				mockGmailService.On("GetLabelsService").Return(mockLabeler)

				if test.isGmailSvcNil {
					mailxSvc.On("GetGmailService", test.userID).Return((*MockGmailService)(nil))
					mailxSvc.On("RecreateGmailService", test.ctx, test.userID).Return(mockGmailService, test.errGmailSvcRecreate)
				}

				if !test.isGmailSvcNil {
					mailxSvc.On("GetGmailService", test.userID).Return(mockGmailService)
				}

				labelsSvc := New(logger, nil, mailxSvc)
				_, err := labelsSvc.GetLabels(test.userID)
				test.assertErr(t, err)
			})
		})
	}

}
