package labels

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-kit/log"
	"github.com/orlandorode97/mailx-google-service"
	"github.com/orlandorode97/mailx-google-service/pkg/google"
	"github.com/orlandorode97/mailx-google-service/pkg/repos"
	"google.golang.org/api/gmail/v1"
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
	svc := s.mailxSvc.GetGmailService(userID)
	if svc == nil || (reflect.ValueOf(svc).Kind() == reflect.Ptr && reflect.ValueOf(svc).IsNil()) {
		return nil
	}
	return svc.GetLabelsService()
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

func (s *service) GetLabels(userID string) ([]*gmail.Label, error) {
	var svc google.Labeler
	var err error
	if svc = s.getLabelService(userID); svc == nil {
		svc, err = s.recreateLabelService(context.Background(), userID)
		if err != nil {
			return nil, err
		}
	}

	labelListCall := svc.List(userID)
	labels, err := labelListCall.Do()

	if err != nil {
		s.logger.Log(
			"message", fmt.Sprintf("error getting labels for user=%s", userID),
			"error", err.Error(),
			"severity", "ERROR",
		)
		return nil, err
	}

	s.logger.Log(
		"message", fmt.Sprintf("get labels for user=%s", userID),
		"severity", "INFO",
	)

	return labels.Labels, nil
}

func (s *service) UpdateLabel() {
}
