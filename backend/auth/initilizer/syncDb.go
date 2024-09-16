package initilizer

import "auth/models"

func SyncDb() {
	DB.AutoMigrate(&models.User{})
}
