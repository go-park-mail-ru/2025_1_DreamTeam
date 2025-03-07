package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
}

func NewEnvironment() *Environment {
	err := godotenv.Load("../../internal/env/.env")
	if err != nil {
		log.Fatalf("Download error .env file: %s", err)
	}

	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_NAME := os.Getenv("DB_NAME")

	if DB_HOST == "" || DB_PORT == "" || DB_USER == "" || DB_PASSWORD == "" || DB_NAME == "" {
		log.Fatal("Missing required environment variables")
	}

	return &Environment{
		DB_HOST:     DB_HOST,
		DB_PORT:     DB_PORT,
		DB_USER:     DB_USER,
		DB_PASSWORD: DB_PASSWORD,
		DB_NAME:     DB_NAME,
	}

}
