package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func Connection(logger *log.Logger) *gorm.DB {

	errEnv := godotenv.Load()

	if errEnv != nil {
		logger.Fatalf("Error loading .env file %v", errEnv.Error())
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	fmt.Printf("host=%s port=%s user=%s password=%s dbname=%s", dbHost, dbPort, dbUser, dbPassword, dbName)

	dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
