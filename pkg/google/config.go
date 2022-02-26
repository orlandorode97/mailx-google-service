package google

import (
	"context"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	oauthv2 "google.golang.org/api/oauth2/v2"
)

type OAuthConfiguration interface {
	Exchange(context.Context, string, ...oauth2.AuthCodeOption) (*oauth2.Token, error)
	AuthCodeURL(string, ...oauth2.AuthCodeOption) string
}

func NewConfig() *oauth2.Config {
	return &oauth2.Config{
		RedirectURL:  viper.GetString("GOOGLE_REDIRECT_URL"),
		ClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
		ClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			gmail.MailGoogleComScope,
			gmail.GmailAddonsCurrentActionComposeScope,
			gmail.GmailAddonsCurrentMessageActionScope,
			gmail.GmailAddonsCurrentMessageMetadataScope,
			gmail.GmailAddonsCurrentMessageReadonlyScope,
			gmail.GmailComposeScope,
			oauthv2.UserinfoProfileScope,
		},
		Endpoint: google.Endpoint,
	}
}
