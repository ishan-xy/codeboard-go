package initializers

import "backend/models"

func InitUserCache() error {
	db := DB
	var users []models.UserPublicProfile
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		models.UserList.Set(user.LeetcodeUsername, &user)
	}

	return nil
}

func StartSyncLoop() {
	go models.SyncLoop()
}
