package initializers

import "backend/models"

func SyncDB() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.UserPublicProfile{})
}