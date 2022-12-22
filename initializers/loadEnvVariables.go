package initializers

import (
	"log"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	projectName := regexp.MustCompile(`^(.*calenderapi)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))
	err := godotenv.Load(string(rootPath) + `/.env`)

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
