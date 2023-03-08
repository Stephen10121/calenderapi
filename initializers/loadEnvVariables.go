package initializers

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	absPath, err2 := filepath.Abs("./.env")
	if err2 != nil {
		log.Fatal("Error finding .env file")
	}

	err := godotenv.Load(absPath)

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
