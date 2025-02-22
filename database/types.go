package database

import "go.mongodb.org/mongo-driver/v2/mongo"

type DifficultyLevel int

const (
	Easy DifficultyLevel = iota
	Medium
	Hard
)

type LeetCodeCacheData struct {
	RealName       string `bson:"real_name"`
	UserAvatar     string `bson:"user_avatar"`
	QuestionCount  int    `bson:"question_count"`
	Ranking        int    `bson:"ranking"`
	TotalQuestions int    `bson:"total_questions"`
}

type Username struct {
	Username string `json:"username"`
}

type Collection[T any] struct {
	*mongo.Collection
}