package initializers

import (
	"os"
	"log"
	"context"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	utils "github.com/ItsMeSamey/go_utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	dsn := os.Getenv("DB")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	
}

func ConnectMongoDB(){
	client, err := mongo.Connect(options.Client().ApplyURI(os.Getenv("MONGO_URI")))

	if err != nil {
		log.Fatalln(utils.WithStack(err))
	}

	// Send a ping to confirm a successful connection
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalln(utils.WithStack(err))
		panic(err)
	}

	log.Println("Pinged your deployment. You successfully connected to MongoDB!")

	// DB = client.Database(os.Getenv("MONGODB_DB"))
	// UserDB = Collection[User]{DB.Collection("users")}
	// ClientDB = Collection[Client]{DB.Collection("clients")}
	// SentMailDB = Collection[Client]{DB.Collection("sent-mails")}
	// log.Println(UserDB.Collection.Name())
}