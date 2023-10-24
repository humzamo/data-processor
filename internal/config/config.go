package config

import (
	"fmt"
	"os"
)

var (
	MongoURI string
)

const defaultMongoURI = "mongodb://username:password@localhost:27017"

func Init() {
	MongoURI = os.Getenv("MONGO_URI")
	if MongoURI == "" {
		MongoURI = defaultMongoURI
	}
	fmt.Println("Mongo URI set to:", MongoURI)
}
