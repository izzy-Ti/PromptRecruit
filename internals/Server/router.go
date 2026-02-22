package server

import (
	//"github.com/gorilla/mux"
	"github.com/gorilla/mux"
	user "github.com/izzy-Ti/PromptRecruit/internals/User"
	"gorm.io/gorm"
)

func Auth(r *mux.Router, db *gorm.DB, jwtSecret string) {
	repo := user.NewUserRepository(db)
	svc := user.NewUserService(repo, jwtSecret)
	h := user.NewAuthHandler(svc)

	userAuth := r.PathPrefix("/user").Subrouter()
	userAuth.HandleFunc("/Login", h.Login).Methods("POST")
	userAuth.HandleFunc("/register", h.Register).Methods("POST")
	userAuth.HandleFunc("/logout", h.Logout).Methods("POST")
	userAuth.HandleFunc("/SendVerifyOTP", h.SendVerifyOTP).Methods("POST")
	userAuth.HandleFunc("/VerifyOTP", h.VerifyOTP).Methods("POST")
	userAuth.HandleFunc("/SendResetOTP", h.SendResetOTP).Methods("POST")
	userAuth.HandleFunc("/ResetPassword", h.ResetPassword).Methods("POST")
	userAuth.HandleFunc("/GoogleAuth", h.GoogleAuth).Methods("POST")
	//userAuth.Handle("/verify", Middleware.IsAuth(svc, http.HandlerFunc(h.VerifyOTP)))
}
