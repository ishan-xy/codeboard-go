package models

import (
	"gorm.io/gorm"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type User struct {
	gorm.Model
	Email            string `gorm:"unique"`
	Name             string
	ProfilePic       string
	LeetcodeUsername string `gorm:"unique"`
	PublicProfile    UserPublicProfile `gorm:"foreignKey:LeetcodeUsername;references:LeetcodeUsername"`
}

type UserPublicProfile struct {
	gorm.Model
	LeetcodeUsername string `gorm:"uniqueIndex"`
	QuestionCount    int
	Ranking          int
	TotalQuestions   int
	Solved           map[string]string `gorm:"type:jsonb"`
}

// Initialize Solved map if nil
func (upp *UserPublicProfile) AfterFind(tx *gorm.DB) (err error) {
	if upp.Solved == nil {
		upp.Solved = make(map[string]string)
	}
	return
}

type Question struct {
	Difficulty    DifficultyLevel `json:"difficulty"`
	Title         string          `json:"title"`
	TitleSlug     string          `json:"titleSlug"`
}

type QuestionList struct {
	Timestamp int64
	Timeout   int64
	SlugMap   map[string]*Question `json:"-"`
	Total     int                  `json:"total"`
	Questions []Question           `json:"questions"`
}

// Populate SlugMap from Questions
func (ql *QuestionList) BuildSlugMap() {
	ql.SlugMap = make(map[string]*Question)
	for i := range ql.Questions {
		ql.SlugMap[ql.Questions[i].TitleSlug] = &ql.Questions[i]
	}
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

func (UserPublicProfile) TableName() string {
	return "user_public_profiles"
}