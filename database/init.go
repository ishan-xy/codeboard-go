package database

import (
	"backend/config"
	"context"
	_ "fmt"
	"log"

	utils "github.com/ItsMeSamey/go_utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var DB *mongo.Database
var UserDB Collection[User]
var QuestionDB Collection[Question]
var PublicUsersDB Collection[UserPublicProfile]


func init() {
	
	client, err := mongo.Connect(options.Client().ApplyURI(config.Cfg.MongoURI))

	if err != nil {
		log.Fatalln(utils.WithStack(err))
	}

	// Send a ping to confirm a successful connection
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalln(utils.WithStack(err))
		panic(err)
	}

	log.Println("Pinged your deployment. You successfully connected to MongoDB!")

	DB = client.Database(config.Cfg.DBName)
	UserDB = Collection[User]{DB.Collection("users")}
	QuestionDB = Collection[Question]{DB.Collection("questions")}
	PublicUsersDB = Collection[UserPublicProfile]{DB.Collection("public-users")}
	log.Println(UserDB.Collection.Name())
}
