package initializers

import (
	"log"
	"path/filepath"

	"gorm.io/driver/sqlite"
	_ "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	databasePath, err2 := filepath.Abs("./test.db")
	if err2 != nil {
		log.Fatal("Failed to find a to database!")
	}

	database, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	DB = database
}
