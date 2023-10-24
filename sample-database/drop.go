package main

import (
	"context"
	"data-processing/internal/config"
	"data-processing/internal/database"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	config.Init()

	clientOptions := options.Client().ApplyURI(config.MongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(database.DatabaseName)
	collection := db.Collection(database.PersonsCollection)

	err = collection.Drop(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("successfully dropped collection:", database.PersonsCollection)
}
