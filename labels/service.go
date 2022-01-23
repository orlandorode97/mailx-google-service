package labels

import (
	kitlog "github.com/go-kit/kit/log"
	"github.com/orlandorode97/mailx-google-service/repos"
	"google.golang.org/api/gmail/v1"
)

type Service interface {
	CreateLabel()
	DeleteLabel()
	GetLabelById()
	GetLabels()
	UpdateLabel()
}

type service struct {
	logger kitlog.Logger
	repo   repos.Repository
	gmail  *gmail.Service
}

func NewService(logger kitlog.Logger, repo repos.Repository, gmail *gmail.Service) Service {
	return &service{
		logger: logger,
		repo:   repo,
		gmail:  gmail,
	}
}

func (s *service) CreateLabel() {

}

func (s *service) DeleteLabel() {

}

func (s *service) GetLabelById() {

}

func (s *service) GetLabels() {

}

func (s *service) UpdateLabel() {

}
