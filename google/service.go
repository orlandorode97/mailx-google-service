package google

import (
	kitlog "github.com/go-kit/kit/log"
	"github.com/orlandorode97/mailx-google-service/repos"
)

type Service interface {
	GetUserProfile()
	CreateLabel()
	DeleteLabel()
	GetLabelById()
	GetLabels()
	UpdateLabel()
}

type service struct {
	logger kitlog.Logger
	repo   repos.Repository
}

func NewService(logger kitlog.Logger, repo repos.Repository) Service {
	return &service{
		logger: logger,
		repo:   repo,
	}
}

func (s *service) GetUserProfile() {

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
