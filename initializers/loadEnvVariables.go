package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load(`/home/pi/golang/src/github.com/stephen10121/calenderapi/.env`)

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
