package main

import (
	"log"
	"net/http"
	"os"

	server "github.com/izzy-Ti/PromptRecruit/internals/Server"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	handler := server.New().PathPrefix("/api/v1").Subrouter()
	server.Auth(handler)

	log.Printf("Listening on " + os.Getenv("PORT"))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), handler))
}
