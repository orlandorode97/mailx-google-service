package mailx

import (
	"context"

	"github.com/go-kit/log"
	"github.com/orlandorode97/mailx-google-service/pkg/repos"
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
	// GetGmailService returns a gmail pointer service by the email of the user.
	GetGmailService(string) *gmail.Service
	// CreateGmailService returns a new gmail service instance.
	CreateGmailService(*oauth2.Token) (*gmail.Service, error)
	// AddGmailServiceByID creates a new entry of a pointer gmail service by google user ID.
	AddGmailServiceByID(string, *gmail.Service)
	// RecreateGmailService returns a new gmail service when a service is not attached to a user
	RecreateGmailService(context.Context, string) (*gmail.Service, error)
}

type service struct {
	logger log.Logger
	//`config` keeps the oauth2 configuration that holds google_client_id, client_secret, and other needed things.
	config *oauth2.Config
	repo   repos.Repository
	//Map that holds google user ID as a key and stores a pointer gmail service.
	gmailSvcs map[string]*gmail.Service
}

func New(logger log.Logger, repo repos.Repository, config *oauth2.Config) Service {
	return &service{
		logger:    logger,
		config:    config,
		repo:      repo,
		gmailSvcs: make(map[string]*gmail.Service),
	}
}

func (s *service) AddGmailServiceByID(ID string, gmailSvc *gmail.Service) {
	s.gmailSvcs[ID] = gmailSvc
}

func (s *service) GetGmailService(ID string) *gmail.Service {
	gmailSvc, ok := s.gmailSvcs[ID]
	if !ok {
		return nil
	}
	return gmailSvc
}

func (s *service) CreateGmailService(token *oauth2.Token) (*gmail.Service, error) {
	ctx := context.Background()
	gmailSvc, err := gmail.NewService(ctx, option.WithTokenSource(s.config.TokenSource(ctx, token)))
	if err != nil {
		s.logger.Log(
			"message", "could not create gmail service",
			"err", err.Error(),
			"severity", "CRITICAL",
		)
		return nil, err
	}
	return gmailSvc, nil
}

func (s *service) RecreateGmailService(ctx context.Context, ID string) (*gmail.Service, error) {
	token, err := s.repo.GetTokenByUserId(ctx, ID)
	if err != nil {
		return nil, err
	}

	svc, err := s.CreateGmailService(&oauth2.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    token.TokenType,
		Expiry:       token.TokenExpiration,
	})
	if err != nil {
		return nil, err
	}
	s.AddGmailServiceByID(ID, svc)
	return svc, nil
}
