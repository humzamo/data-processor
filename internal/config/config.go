package config

import (
	"fmt"
	"os"
)

var (
	MongoURI string
	Port     string
)

const (
	defaultMongoURI = "mongodb://username:password@localhost:27017"
	defaultPort     = "8080"
)

// Init initilises the config for the microservice
func Init() {
	MongoURI = os.Getenv("MONGO_URI")
	if MongoURI == "" {
		MongoURI = defaultMongoURI
	}
	fmt.Println("Mongo URI set to:", MongoURI)

	Port = os.Getenv("PORT")
	if Port == "" {
		Port = defaultPort
	}
	fmt.Println("Port set to:", Port)
}
