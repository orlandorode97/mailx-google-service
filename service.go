package mailx

import (
	"context"

	"github.com/go-kit/log"
	"github.com/orlandorode97/mailx-google-service/pkg/google"
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

type Creator interface {
	// CreateGmailService returns a new gmail service instance.
	CreateGmailService(*oauth2.Token) (google.Service, error)
	// RecreateGmailService returns a new gmail service when a service is not attached to a user
	RecreateGmailService(context.Context, string) (google.Service, error)
}

type Getter interface {
	// GetGmailService returns a pointer of gmail service.
	GetGmailService(string) google.Service
}

type Setter interface {
	// AddGmailServiceByID creates a new entry of a pointer gmail service by google user ID.
	AddGmailServiceByID(string, google.Service) google.Service
}

type Service interface {
	Creator
	Getter
	Setter
}

type service struct {
	logger log.Logger
	//`config` keeps the oauth2 configuration that holds google_client_id, client_secret, and other needed things.
	config *oauth2.Config
	repo   repos.TokenRepository
	//Map that holds google user ID as a key and stores a pointer gmail service.
	gmailSvcs map[string]google.Service
}

func New(logger log.Logger, repo repos.TokenRepository, config *oauth2.Config) Service {
	return &service{
		logger:    logger,
		config:    config,
		repo:      repo,
		gmailSvcs: make(map[string]google.Service),
	}
}

func (s *service) AddGmailServiceByID(userID string, gmailSvc google.Service) google.Service {
	s.gmailSvcs[userID] = gmailSvc
	return gmailSvc
}

func (s *service) GetGmailService(userID string) google.Service {
	if gmailSvc, ok := s.gmailSvcs[userID]; ok {
		return gmailSvc
	}

	return nil
}

func (s *service) CreateGmailService(token *oauth2.Token) (google.Service, error) {
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

	svc := s.hydrateServices(gmailSvc)

	return svc, nil
}

func (s *service) RecreateGmailService(ctx context.Context, userID string) (google.Service, error) {
	token, err := s.repo.GetTokenByUserId(ctx, userID)
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

	return s.AddGmailServiceByID(userID, svc), nil
}

func (s *service) hydrateServices(svc *gmail.Service) google.Service {
	return &google.GmailService{
		Users:    svc.Users,
		Labels:   google.NewLabelsService(svc.Users.Labels),
		Messages: google.NewMessagesService(svc.Users.Messages),
		Drafts:   &google.DraftsService{},
		History:  &google.HistoryService{},
		Settings: &google.SettingsService{},
		Threads:  &google.ThreadsService{},
	}
}
