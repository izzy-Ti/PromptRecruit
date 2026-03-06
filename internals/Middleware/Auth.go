package Middleware

import (
	"context"
	"errors"
	"net/http"

	models "github.com/izzy-Ti/PromptRecruit/internals/Models"
	//user "github.com/izzy-Ti/PromptRecruit/internals/User"
	"github.com/izzy-Ti/PromptRecruit/internals/Utils"
)

type TokenValidator interface {
	ValidateToken(token string) (*models.User, error)
}

func IsAuth(svc TokenValidator, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			Utils.WriteJson(w, http.StatusUnauthorized, errors.New("Unauthorized please login"))
			return
		}

		user, err := svc.ValidateToken(token.Value)
		if err != nil {
			Utils.WriteJson(w, http.StatusUnauthorized, err)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
