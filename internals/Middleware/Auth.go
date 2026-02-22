package Middleware

import (
	"context"
	"errors"
	"net/http"

	user "github.com/izzy-Ti/PromptRecruit/internals/User"
	"github.com/izzy-Ti/PromptRecruit/internals/Utils"
)

func IsAuth(svc *user.UserService, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			Utils.WriteJson(w, http.StatusUnauthorized, errors.New("Unauthorized please login"))
			return
		}
		user, err := svc.ValidateToken(token.Value)
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
