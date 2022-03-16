package mailx

import (
	"context"

	"github.com/go-kit/log"
	"github.com/orlandorode97/mailx-google-service/pkg/repos"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// Type of service to return by GetGmailService.
const (
	// Returns gmailSvc.Users service.
	UsersSvc = iota + 1
	// Returns gmailSvc.Users.Labels service.
	LabelSvc
	// Returns gmailSvc.Users.Drafts service.
	DraftsSvc
	// Returns gmailSvc.Users.History service.
	HistorySvc
	// Returns gmailSvc.Users.Messages service.
	MessagesSvc
	// Returns gmailSvc.Users.Settings service.
	SettingsSvc
	// Returns gmailSvc.Users.Threads service.
	ThreadsSvc
	// Returns the whole gmailSvc
	GmailSvc
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
	CreateGmailService(*oauth2.Token) (*gmail.Service, error)
	// RecreateGmailService returns a new gmail service when a service is not attached to a user
	RecreateGmailService(context.Context, string) (*gmail.Service, error)
}

type Getter interface {
	// GetGmailService returns a interface of any gmail service such as Labels, Users, Drafts etc of the user based on typeSvc parameter.
	GetGmailService(string, int) interface{}
}

type Setter interface {
	// AddGmailServiceByID creates a new entry of a pointer gmail service by google user ID.
	AddGmailServiceByID(string, *gmail.Service) *gmail.Service
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
	gmailSvcs map[string]*gmail.Service
}

func New(logger log.Logger, repo repos.TokenRepository, config *oauth2.Config) Service {
	return &service{
		logger:    logger,
		config:    config,
		repo:      repo,
		gmailSvcs: make(map[string]*gmail.Service),
	}
}

func (s *service) AddGmailServiceByID(ID string, gmailSvc *gmail.Service) *gmail.Service {
	s.gmailSvcs[ID] = gmailSvc
	return gmailSvc
}

func (s *service) GetGmailService(ID string, typeSvc int) interface{} {
	gmailSvc, ok := s.gmailSvcs[ID]
	if !ok {
		return nil
	}
	switch typeSvc {
	case UsersSvc:
		return gmailSvc.Users
	case LabelSvc:
		return gmailSvc.Users.Labels
	case DraftsSvc:
		return gmailSvc.Users.Drafts
	case HistorySvc:
		return gmailSvc.Users.History
	case MessagesSvc:
		return gmailSvc.Users.Messages
	case SettingsSvc:
		return gmailSvc.Users.Settings
	case ThreadsSvc:
		return gmailSvc.Users.Threads
	default:
		return gmailSvc
	}
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

	return s.AddGmailServiceByID(ID, svc), nil
}
