package google

import (
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

func NewConfig() (*oauth2.Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

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
		},
		Endpoint: google.Endpoint,
	}, nil
}
