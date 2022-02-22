package middlewares

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/orlandorode97/mailx-google-service/auth"
	"github.com/orlandorode97/mailx-google-service/pkg/models"
	"github.com/spf13/viper"
)

type MailxClaims struct {
	ID string
	jwt.StandardClaims
}

type contextMailxKey string

var (
	InvalidAuthKey contextMailxKey = "InvalidAuth"
)

const MailxJWTHeader string = "mailx_google_auth"

func JWTKeyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(viper.GetString("JWT_SIGNING_KEY")), nil
}

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		cookie, err := r.Cookie("mailx_google_auth")
		if err != nil {
			ctx = context.WithValue(r.Context(), InvalidAuthKey, models.ErrInvalidCookie{})
			next.ServeHTTP(rw, r.WithContext(ctx))
		}

		_, err = jwt.ParseWithClaims(cookie.Value, &auth.MailxClams{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("JWT_SIGNING_KEY")), nil
		})

		if err != nil {
			e, ok := err.(*jwt.ValidationError)
			if ok {
				switch {
				case e.Errors&jwt.ValidationErrorSignatureInvalid != 0:
					ctx = context.WithValue(r.Context(), InvalidAuthKey, models.ErrInvalidSignature{})
				case e.Errors&jwt.ValidationErrorMalformed != 0:
					ctx = context.WithValue(r.Context(), InvalidAuthKey, models.ErrMalformedToken{})
				case e.Errors&jwt.ValidationErrorNotValidYet != 0:
					ctx = context.WithValue(r.Context(), InvalidAuthKey, models.ErrInactiveToken{})
				case e.Errors&jwt.ValidationErrorExpired != 0:
					ctx = context.WithValue(r.Context(), InvalidAuthKey, models.ErrExpiredToken{})
				case e.Inner != nil:
					ctx = context.WithValue(r.Context(), InvalidAuthKey, e.Inner)

				}
			}
		}

		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
