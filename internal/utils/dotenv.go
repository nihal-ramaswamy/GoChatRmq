package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetDotEnvVariable(key string) string {
	err := godotenv.Load()

	if nil != err {
		log.Fatalf("Error loading .env file")
	}

	value := os.Getenv(key)

	return value
}
