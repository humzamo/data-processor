package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"data-processing/internal/config"
	"data-processing/internal/models"
)

var DB *mongo.Database

const (
	databaseName      = "sampleDatabase"
	PersonsCollection = "persons"
	BatchSize         = 10
	DefaultTimeout    = time.Second * 10
)

func InitDB() {
	clientOptions := options.Client().ApplyURI(config.MongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set the database and collection
	DB = client.Database(databaseName)
	collection := DB.Collection(PersonsCollection)

	// Create the indexes
	_, err = collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}
	fmt.Println("Indexes created successfully.")

	// Add sample data if collection is empty
	count, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		log.Fatalf("Failed to count documents in collection: %v", err)
	}
	fmt.Println("Initial document count:", count)

	if count == 0 {
		_, err = collection.InsertMany(ctx, sampleData)
		if err != nil {
			log.Fatalf("Failed to add sample data to collection: %v", err)
		}
		fmt.Println("Sample data added to collection successfully.")
	}

}

type ProcessedItems struct {
	Count int `json:"processedItems"`
}

func GetProcessedItemCount() (int, error) {
	collection := DB.Collection(PersonsCollection)

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.M{"processed": true})
	if err != nil {
		return -1, errors.Wrap(err, "unable to get count of processed items")
	}

	return int(count), nil
}

// GetUnprocessedPersons retrieves a batch of unprocessed persons from the database.
func GetUnprocessedPersons(batchSize int) ([]models.Person, error) {
	var persons []models.Person

	collection := DB.Collection(PersonsCollection)

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "_id", Value: 1}})
	findOptions.SetLimit(int64(batchSize))

	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	filter := bson.M{"processed": false}

	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &persons); err != nil {
		return nil, err
	}

	return persons, nil
}