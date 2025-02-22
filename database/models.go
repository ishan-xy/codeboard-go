package database

import (
	cmap "github.com/orcaman/concurrent-map/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID     `bson:"_id,omitempty"`
	Email         string            `json:"email" bson:"email"`
	Name          string            `json:"name" bson:"name"`
	PublicProfileID primitive.ObjectID   `json:"public_profile" bson:"public_profile"`
}

type UserPublicProfile struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	LeetcodeUsername string        `bson:"leetcode_username"`
	ProfilePic       string        `json:"profile_pic" bson:"profile_pic"`
	QuestionCount    int           `bson:"question_count"`
	Ranking          int           `bson:"ranking"`
	TotalQuestions   int           `bson:"total_questions"`
	Solved           []Question    `bson:"solved"`
}

type Question struct {
	ID         primitive.ObjectID   `bson:"_id,omitempty"`
	Difficulty DifficultyLevel `json:"difficulty" bson:"difficulty"`
	Title      string          `json:"title" bson:"title"`
	TitleSlug  string          `json:"titleSlug" bson:"title_slug"`
}

type QuestionList struct {
	ID        primitive.ObjectID        `bson:"_id,omitempty"`
	Timestamp int64                `bson:"timestamp"`
	Timeout   int64                `bson:"timeout"`
	SlugMap   map[string]*Question `json:"-" bson:"-"`
	Total     int                  `json:"total" bson:"total"`
	Questions []Question           `json:"questions" bson:"questions"`
}

// BuildSlugMap initializes the slug map from questions
func (ql *QuestionList) BuildSlugMap() {
	ql.SlugMap = make(map[string]*Question)
	for i := range ql.Questions {
		ql.SlugMap[ql.Questions[i].TitleSlug] = &ql.Questions[i]
	}
}

type PublicUserCache struct {
	PublicUsers cmap.ConcurrentMap[string, LeetCodeCacheData]
}

var UserCache = &PublicUserCache{
	PublicUsers: cmap.New[LeetCodeCacheData](),
}


