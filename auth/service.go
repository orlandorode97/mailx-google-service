package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	kitlog "github.com/go-kit/kit/log"
	"github.com/orlandorode97/mailx-google-service"
	"github.com/orlandorode97/mailx-google-service/google"
	"github.com/orlandorode97/mailx-google-service/repos"
	"golang.org/x/oauth2"
)

type Service interface {
	GetOauthUrl(context.Context) (string, error)
	GenerateOauthToken(context.Context, string) (*oauth2.Token, error)
	GenerateHttpClient(context.Context, *oauth2.Token) *http.Client
	ConfigGmailService(context.Context, *http.Client) error
	CreateUser(context.Context, *oauth2.Token) error
}

type service struct {
	logger       kitlog.Logger
	db           repos.Repository
	mailxService mailx.Service
}

func NewService(logger kitlog.Logger, db repos.Repository, mailx mailx.Service) Service {
	return &service{
		logger:       logger,
		db:           db,
		mailxService: mailx,
	}
}

func (s *service) GetOauthUrl(_ context.Context) (string, error) {
	state, err := randomState()
	if err != nil {
		return "", err
	}

	url := s.mailxService.GetOauthConfig().AuthCodeURL(state, oauth2.AccessTypeOffline)
	if url == "" {
		return "", mailx.ErrAuthUrl{}
	}
	return url, nil
}

func (s *service) GenerateOauthToken(_ context.Context, code string) (*oauth2.Token, error) {
	token, err := s.mailxService.GetOauthConfig().Exchange(context.Background(), code)
	if err != nil {
		s.logger.Log(
			"message", "could not create oauth2 token",
			"err", err.Error(),
			"severity", "CRITICAL",
		)
		return nil, err
	}
	return token, nil
}

func (s *service) GenerateHttpClient(_ context.Context, token *oauth2.Token) *http.Client {
	return google.NewClient(s.mailxService.GetOauthConfig(), token)
}

func (s *service) ConfigGmailService(_ context.Context, client *http.Client) error {
	return s.mailxService.ConfigGmailService(client)
}

func (s *service) CreateUser(ctx context.Context, token *oauth2.Token) error {
	response, err := s.mailxService.GetHttpClient().Get(mailx.UserInfoUrl + token.AccessToken)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	var user *repos.User

	if err = json.NewDecoder(response.Body).Decode(&user); err != nil {
		return err
	}

	_, err = s.db.GetUserByID(ctx, user.ID)
	if err == sql.ErrNoRows {
		s.logger.Log(
			"message", fmt.Sprintf("creating user %s with ID %s", user.GivenName, user.ID),
			"severity", "INFO",
		)
		if err = s.db.CreateUser(ctx, user); err != nil {
			s.logger.Log(
				"message", fmt.Sprintf("error creating user %s with ID %s", user.GivenName, user.ID),
				"error", err.Error(),
				"severity", "ERROR",
			)
			return err
		}
		return nil
	}

	if err != nil {
		s.logger.Log(
			"message", fmt.Sprintf("error by getting user with ID %s", user.ID),
			"error", err.Error(),
			"severity", "ERROR",
		)
		return err
	}

	return nil
}

func randomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
