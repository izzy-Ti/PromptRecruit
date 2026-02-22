package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/izzy-Ti/PromptRecruit/internals/Middleware"
	user "github.com/izzy-Ti/PromptRecruit/internals/User"
	"gorm.io/gorm"
)

func Auth(r *mux.Router, db *gorm.DB, jwtSecret string) {
	repo := user.NewUserRepository(db)
	svc := user.NewUserService(repo, jwtSecret)
	h := user.NewAuthHandler(svc)

	userAuth := r.PathPrefix("/user").Subrouter()
	userAuth.HandleFunc("/Login", h.Login).Methods("POST")
	userAuth.Handle("/verify", Middleware.IsAuth(svc, http.HandlerFunc(h.VerifyOTP)))
}
