package config

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUsername string
	DBPassword string
	DBName     string
}

func Connection() *gorm.DB {

	fg := Config{
		DBHost:     "postgres",
		DBPort:     "5432",
		DBUsername: "postgres",
		DBPassword: "postgres",
		DBName:     "dynamic_segment_service_db",
	}

	fmt.Printf("host=%s port=%s user=%s password=%s dbname=%s", fg.DBHost, fg.DBPort, fg.DBUsername, fg.DBPassword, fg.DBName)

	dbUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", fg.DBHost, fg.DBPort, fg.DBUsername, fg.DBPassword, fg.DBName)
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	return db
}
