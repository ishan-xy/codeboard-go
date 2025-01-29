package models

import (
	cmap "github.com/orcaman/concurrent-map/v2"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email            string `gorm:"unique"`
	Name             string
	ProfilePic       string
	LeetcodeUsername string `gorm:"unique"`  // Unique identifier for public profile
	PublicProfile    UserPublicProfile `gorm:"foreignKey:LeetcodeUsername;references:LeetcodeUsername"`
}

type UserPublicProfile struct {
	gorm.Model
	LeetcodeUsername string `gorm:"uniqueIndex"`  // Foreign key connecting to User
	QuestionCount    int
	Ranking          int
	TotalQuestions   int
}

// models/users.go
type LeetCodeCacheData struct {
	RealName       string
	UserAvatar     string
	QuestionCount  int
	Ranking        int
	TotalQuestions int
}

type PublicUserCache struct {
	PublicUsers cmap.ConcurrentMap[string, LeetCodeCacheData]
}

var UserCache = &PublicUserCache{
	PublicUsers: cmap.New[LeetCodeCacheData](),
}

// Belongs-to relationship: UserPublicProfile belongs to User through LeetcodeUsername
func (UserPublicProfile) TableName() string {
	return "user_public_profiles"
}