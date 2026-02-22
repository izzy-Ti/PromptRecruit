package main

import (
	"log"
	"net/http"
	"os"
	"sync"

	db "github.com/izzy-Ti/PromptRecruit/internals/Db"
	server "github.com/izzy-Ti/PromptRecruit/internals/Server"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		db.Connect()
		db.Migrate()
	}()

	handler := server.New().PathPrefix("/api/v1").Subrouter()
	server.Auth(handler, db.DB, os.Getenv("JWT_SECRET"))

	wg.Wait()

	log.Print("Listening on " + os.Getenv("PORT"))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), handler))
}
