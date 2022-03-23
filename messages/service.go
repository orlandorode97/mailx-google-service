package messages

import (
	"context"

	"github.com/go-kit/log"
	"github.com/orlandorode97/mailx-google-service"
	"github.com/orlandorode97/mailx-google-service/pkg/google"
	"github.com/orlandorode97/mailx-google-service/pkg/repos"
	"google.golang.org/api/gmail/v1"
)

type Service interface {
	GetMessages(context.Context, string) ([]*gmail.Message, error)
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

func (s *service) getMessageService(userID string) google.Messenger {
	svc := s.mailxSvc.GetGmailService(userID)
	if svc == nil {
		return nil
	}

	return svc.GetMessagesService()
}

func (s *service) recreateMessageService(ctx context.Context, userID string) (google.Messenger, error) {
	_, err := s.mailxSvc.RecreateGmailService(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.getMessageService(userID), nil
}

func (s *service) GetMessages(ctx context.Context, userID string) ([]*gmail.Message, error) {
	var svc google.Messenger
	var err error
	if svc = s.getMessageService(userID); svc == nil {
		svc, err = s.recreateMessageService(context.Background(), userID)
		if err != nil {
			return nil, err
		}
	}

	messages, err := svc.List(userID).Do()
	if err != nil {
		return nil, err
	}

	return messages.Messages, nil
}
