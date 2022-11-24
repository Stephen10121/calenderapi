package initializers

import (
	"github.com/stephen10121/calenderapi/models"
)

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Group{})
	DB.AutoMigrate(&models.Job{})
}
