package db

import (
	"log"
	"os"

	models "github.com/izzy-Ti/PromptRecruit/internals/Models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	log.Println("db connected")
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Discard})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}
func Migrate() {
	err := DB.AutoMigrate(&models.User{}, &models.Cvs{}, &models.Jobs{})
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Database migrated")
}
