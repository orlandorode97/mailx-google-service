package labels

import (
	kitlog "github.com/go-kit/kit/log"
	"github.com/orlandorode97/mailx-google-service"
	"github.com/orlandorode97/mailx-google-service/pkg/repos"
)

const (
	UserInfoUrl = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

type Service interface {
	CreateLabel()
	DeleteLabel()
	GetLabelById()
	GetLabels()
	UpdateLabel()
}

type service struct {
	logger   kitlog.Logger
	repo     repos.Repository
	mailxSvc mailx.Service
}

func NewService(logger kitlog.Logger, repo repos.Repository, mailx mailx.Service) Service {
	return &service{
		logger:   logger,
		repo:     repo,
		mailxSvc: mailx,
	}
}

func (s *service) CreateLabel() {
}

func (s *service) DeleteLabel() {

}

func (s *service) GetLabelById() {

}

func (s *service) GetLabels() {
	// response, err := s.gmail.Users.Labels.List("me").Do()
	// if err != nil {
	// 	return
	// }
	// if len(response.Labels) == 0 {
	// 	return
	// }
	// fmt.Println("Labels")
	// for _, label := range response.Labels {
	// 	fmt.Printf("- %s \n", label.Name)
	// }
}

func (s *service) UpdateLabel() {

}
