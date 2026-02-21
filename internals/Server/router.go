package server

import (
	"github.com/gorilla/mux"
	user "github.com/izzy-Ti/PromptRecruit/internals/User"
)

func Auth(r *mux.Router) {
	userAuth := r.PathPrefix("/user").Subrouter()
	userAuth.HandleFunc("/Login", user.Login).Methods("POST")
}
