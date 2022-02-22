package auth

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/golang-jwt/jwt/v4"
	"github.com/orlandorode97/mailx-google-service"
	"github.com/orlandorode97/mailx-google-service/pkg/models"
	"github.com/orlandorode97/mailx-google-service/pkg/repos"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

type MailxClams struct {
	ID string
	jwt.StandardClaims
}

type Service interface {
	// Creates the oath authorization URL
	GetOauthUrl(context.Context) (string, error)
	// Generates the token access after a successful sign in
	GenerateOauthToken(context.Context, string) (*oauth2.Token, error)
	// Configuration of a gmail service for a current user
	ConfigGmailServiceUser(context.Context, string) (*models.User, error)
	// CreateJWT generates a json web token for mailx-google-service authentication.
	CreateJWT(context.Context, *models.User) (string, error)
}

type service struct {
	logger       kitlog.Logger
	db           repos.Repository
	client       *http.Client
	mailxService mailx.Service
}

// NewService creates a new Auth Service.
func NewService(logger kitlog.Logger, db repos.Repository, mailx mailx.Service) Service {
	return &service{
		logger:       logger,
		db:           db,
		client:       http.DefaultClient,
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
		return "", models.ErrAuthUrl{}
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

func (s *service) ConfigGmailServiceUser(ctx context.Context, code string) (*models.User, error) {
	token, err := s.GenerateOauthToken(ctx, code)
	if err != nil {
		return nil, err
	}

	svc, err := s.mailxService.CreateGmailService(token)
	if err != nil {
		return nil, err
	}

	user, err := s.createUser(ctx, token)
	if err != nil {
		return nil, err
	}

	if err := s.saveAccessToken(ctx, user.ID, token); err != nil {
		return nil, err
	}
	s.mailxService.AddGmailServiceByID(user.ID, svc)
	return user, nil
}

func (s *service) CreateJWT(_ context.Context, user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MailxClams{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Second * 1).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
		},
	})
	return token.SignedString([]byte(viper.GetString("JWT_SIGNING_KEY")))
}

func (s *service) createUser(ctx context.Context, token *oauth2.Token) (*models.User, error) {
	response, err := s.client.Get(mailx.UserInfoUrl + token.AccessToken)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var user *models.User

	if err = json.NewDecoder(response.Body).Decode(&user); err != nil {
		return nil, err
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
			return nil, err
		}
		return user, nil
	}

	if err != nil {
		s.logger.Log(
			"message", fmt.Sprintf("error by getting user with ID %s", user.ID),
			"error", err.Error(),
			"severity", "ERROR",
		)
		return nil, err
	}

	s.logger.Log(
		"message", fmt.Sprintf("user with ID %s already exists, skipping creation.", user.ID),
		"severity", "INFO",
	)

	return user, nil
}

func (s *service) saveAccessToken(ctx context.Context, ID string, token *oauth2.Token) error {
	_, err := s.db.GetTokenByUserId(ctx, ID)
	if err == sql.ErrNoRows {
		s.logger.Log(
			"message", fmt.Sprintf("saving token for the user with ID %s", ID),
			"severity", "INFO",
		)
		if err = s.db.SaveAccessToken(ctx, ID, token); err != nil {
			s.logger.Log(
				"message", fmt.Sprintf("error saving token for the userwith ID %s", ID),
				"error", err.Error(),
				"severity", "ERROR",
			)

			return err
		}

		return nil
	}

	if err != nil {
		s.logger.Log(
			"message", fmt.Sprintf("error by token for the user with ID %s", ID),
			"error", err.Error(),
			"severity", "ERROR",
		)
		return err
	}

	if err := s.db.UpdateAccessToken(ctx, ID, token); err != nil {
		return err
	}

	return nil
}

// random state creates a random string to be placed in the authorization URL.
func randomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
