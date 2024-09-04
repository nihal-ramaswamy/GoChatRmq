package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func GetDotEnvVariable(key string) string {
	envFile := os.Getenv("ENVIRONMENT")
	if envFile == "" {
		envFile = dir(".env")
	} else {
		envFile = dir(".env." + envFile)
	}

	err := godotenv.Load(envFile)

	if nil != err {
		log.Fatalf("Error loading .env file: %s", err)
	}

	value := os.Getenv(key)

	return value
}

func dir(envFile string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return filepath.Join(currentDir, envFile)
}
