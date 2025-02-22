package config

import (
	"time"
	"log"

	utils "github.com/ItsMeSamey/go_utils"
)

type Config struct {
	Port        string
	MongoURI    string
	DBName      string
	Secret   string
	JWTExpiration time.Duration
	CookieName 	string

}

var Cfg *Config

func init() {
	loadEnv()
	var err error
	Cfg, err = loadConfig()
	if err != nil {
		log.Fatal(utils.WithStack(err))
	}
	log.Println("Configuration loaded successfully:", Cfg)
}

func loadConfig() (*Config, error) {

	return &Config{
		Port:        Getenv("PORT"),
		MongoURI:    Getenv("MONGO_URI"),
		DBName:      Getenv("MONGODB_DB"),
		Secret:   Getenv("SECRET"),
		JWTExpiration: time.Hour * 24,
		CookieName:  "sessionID",
	}, nil
}
