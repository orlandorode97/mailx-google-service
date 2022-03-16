package labels

import (
	"context"

	"github.com/go-kit/log"
	"github.com/orlandorode97/mailx-google-service"
	"github.com/orlandorode97/mailx-google-service/pkg/google"
	"github.com/orlandorode97/mailx-google-service/pkg/repos"
	"google.golang.org/api/gmail/v1"
)

const (
	UserInfoUrl = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

type Service interface {
	CreateLabel()
	DeleteLabel()
	GetLabelById()
	GetLabels(string) ([]*gmail.Label, error)
	UpdateLabel()
}

type service struct {
	logger   log.Logger
	repo     repos.Repository
	mailxSvc mailx.Service
}

func New(logger log.Logger, repo repos.Repository, mailx mailx.Service) Service {
	return &service{
		logger:   logger,
		repo:     repo,
		mailxSvc: mailx,
	}
}

func (s *service) getLabelService(userID string) google.Labeler {
	svc := s.mailxSvc.GetGmailService(userID, mailx.LabelSvc)
	if svc == nil {
		return nil
	}
	labelsSvc, ok := svc.(google.Labeler)
	if !ok {
		return nil
	}
	return labelsSvc
}

func (s *service) recreateLabelService(ctx context.Context, userID string) (google.Labeler, error) {
	_, err := s.mailxSvc.RecreateGmailService(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.getLabelService(userID), nil
}

func (s *service) CreateLabel() {
}

func (s *service) DeleteLabel() {

}

func (s *service) GetLabelById() {

}

func (s *service) GetLabels(userId string) ([]*gmail.Label, error) {
	var svc google.Labeler
	var err error
	if svc = s.getLabelService(userId); svc == nil {
		svc, err = s.recreateLabelService(context.Background(), userId)
		if err != nil {
			return nil, err
		}
	}

	labels, err := s.doListLabel(svc.List(userId))

	if err != nil {
		return nil, err
	}

	return labels, nil
}

func (s *service) UpdateLabel() {

}

// doListLabel executes the Do request of the List method of gmail.UserLabelsService
func (s *service) doListLabel(labeler google.LabelerClientList) ([]*gmail.Label, error) {
	labels, err := labeler.Do()
	if err != nil {
		return nil, err
	}
	return labels.Labels, nil
}

// doLabel executes the Do request of the rest of the methods of gmail.UserLabelsService
func (s *service) doLabel(labeler google.LabelerClient) (*gmail.Label, error) {
	label, err := labeler.Do()
	if err != nil {
		return nil, err
	}
	return label, nil
}
