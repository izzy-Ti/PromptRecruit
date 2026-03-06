package server

import (
	//"github.com/gorilla/mux"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/izzy-Ti/PromptRecruit/internals/Middleware"
	user "github.com/izzy-Ti/PromptRecruit/internals/User"
	"github.com/izzy-Ti/PromptRecruit/internals/cvs"
	"gorm.io/gorm"
)

func SetupRoutes(r *mux.Router, db *gorm.DB, jwtSecret string) {

	userRepo := user.NewUserRepository(db)
	userSvc := user.NewUserService(userRepo, jwtSecret)

	Crepo := cvs.NewCvRepo(db)
	Csvc := cvs.NewUserService(Crepo)

	Auth(r, userSvc)
	CV(r, db, userSvc, Csvc)
}

func Auth(r *mux.Router, svc *user.UserService) {

	h := user.NewAuthHandler(svc)

	userAuth := r.PathPrefix("/user").Subrouter()
	userAuth.HandleFunc("/login", h.Login).Methods("POST")
	userAuth.HandleFunc("/register", h.Register).Methods("POST")
	userAuth.HandleFunc("/logout", h.Logout).Methods("POST")
	userAuth.HandleFunc("/SendVerifyOTP", h.SendVerifyOTP).Methods("POST")
	userAuth.HandleFunc("/VerifyOTP", h.VerifyOTP).Methods("POST")
	userAuth.HandleFunc("/SendResetOTP", h.SendResetOTP).Methods("POST")
	userAuth.HandleFunc("/ResetPassword", h.ResetPassword).Methods("POST")
	userAuth.HandleFunc("/GoogleAuth", h.GoogleAuth).Methods("POST")
	//userAuth.Handle("/verify", Middleware.IsAuth(svc, http.HandlerFunc(h.VerifyOTP)))
}

func CV(r *mux.Router, db *gorm.DB, userSvc *user.UserService, Cvs *cvs.CVservice) {

	h := cvs.NewCVHnadler(Cvs)

	cvRoute := r.PathPrefix("/cv").Subrouter()
	cvRoute.Handle("/Uploadcv", Middleware.IsAuth(userSvc, http.HandlerFunc(h.CVUploader))).Methods("POST")
	cvRoute.Handle("/apply/{jobId}", Middleware.IsAuth(userSvc, http.HandlerFunc(h.Application))).Methods("POST")
	cvRoute.Handle("/postjob", Middleware.IsAuth(userSvc, http.HandlerFunc(h.JobPost))).Methods("POST")
}
