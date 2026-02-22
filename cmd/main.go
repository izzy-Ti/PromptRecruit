package main

import (
	"log"
	"net/http"
	"os"

	db "github.com/izzy-Ti/PromptRecruit/internals/Db"
	server "github.com/izzy-Ti/PromptRecruit/internals/Server"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	db.Connect()
	db.Migrate()
	handler := server.New().PathPrefix("/api/v1").Subrouter()
	server.Auth(handler, db.DB, os.Getenv("JWT_SECRET"))

	log.Print("Listening on " + os.Getenv("PORT"))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), handler))
}
