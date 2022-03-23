package users

import (
	"context"
	"fmt"

	"github.com/go-kit/log"
	"github.com/orlandorode97/mailx-google-service"
	"github.com/orlandorode97/mailx-google-service/pkg/models"
	"github.com/orlandorode97/mailx-google-service/pkg/repos"
)

type Service interface {
	GetUserByID(string) (*models.User, error)
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

func (s *service) GetUserByID(ID string) (*models.User, error) {
	user, err := s.repo.GetUserByID(context.Background(), ID)
	if err != nil {
		s.logger.Log(
			"message", fmt.Sprintf("error getting user with ID %s ", ID),
			"error", err.Error(),
			"severity", "ERROR",
		)
		return nil, err
	}
	s.logger.Log(
		"message", fmt.Sprintf("getting user with ID %s", ID),
		"severity", "INFO",
	)
	return user, nil
}
