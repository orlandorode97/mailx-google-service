package mailx

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

const (
	UserInfoUrl = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

/*
	Service defines global methods to handle a gmail client and works with the oauth2 configuration.
	This Service is used by multiple services such as:
		- Auth Service
		- Label Service
*/
type Service interface {
	// GetOauthConfig returns the oauth pointer configuration.
	GetOauthConfig() *oauth2.Config
	// GetHttpClient returns an http pointer client.
	GetHttpClient() *http.Client
	// GetGmailService returns a gmail pointer service.
	GetGmailService() *gmail.Service
	// ConfigGmailService sets a new gmail r
	ConfigGmailService(*http.Client) error
}

type service struct {
	logger   log.Logger
	config   *oauth2.Config
	client   *http.Client
	gmailSvc *gmail.Service
}

func NewService(logger log.Logger, config *oauth2.Config) Service {
	return &service{
		logger:   logger,
		config:   config,
		client:   nil,
		gmailSvc: nil,
	}
}

func (s *service) AddService(gmailSvc *gmail.Service) {
	s.gmailSvc = gmailSvc
}

func (s *service) AddClient(client *http.Client) {
	s.client = client
}

func (s *service) GetOauthConfig() *oauth2.Config {
	return s.config
}

func (s *service) GetHttpClient() *http.Client {
	return s.client
}

func (s *service) GetGmailService() *gmail.Service {
	return s.gmailSvc
}

func (s *service) ConfigGmailService(client *http.Client) error {
	gmailSvc, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		s.logger.Log(
			"message", "could not create gmail service",
			"err", err.Error(),
			"severity", "CRITICAL",
		)
		return err
	}
	s.AddClient(client)
	s.AddService(gmailSvc)
	return nil
}
