package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"data-processing/internal/database"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	processingInProgress bool
	processingFinished   bool
)

const (
	batchSize = 3
)

func StartProcessing(w http.ResponseWriter, r *http.Request) {
	fmt.Println("call to StartProcessing handler")

	if processingInProgress {
		http.Error(w, "Processing is already in progress", http.StatusTooManyRequests)
		return
	}

	// Set processing in progress and start processing
	processingInProgress = true

	go func() {
		for {
			// Get a batch of unprocessed persons from the database
			persons, err := database.GetUnprocessedPersons(batchSize)
			if err != nil {
				http.Error(w, "Error fetching unprocessed persons", http.StatusInternalServerError)
				processingInProgress = false
				return
			}

			if len(persons) == 0 {
				// No more unprocessed persons, stop processing
				fmt.Println("Processing has finshed!")
				processingFinished = true
				break
			}

			ctx, cancel := context.WithTimeout(context.Background(), database.DefaultTimeout)
			defer cancel()

			collection := database.DB.Collection(database.PersonsCollection)

			// Process the batch of persons
			for _, person := range persons {
				person.MiddleNames = ExtractMiddleName(person.Name)

				// Update the "processed" flag to mark it as processed
				person.Processed = true

				updateResult, err := collection.ReplaceOne(ctx, bson.M{"_id": person.ID}, person)
				if err != nil {
					log.Printf("Error updating person: %v", err)
				} else {
					log.Printf("Updated %v documents", updateResult.ModifiedCount)
				}
			}

			// Artifical sleep added for testing purposes
			fmt.Println("sleeping...")
			time.Sleep(time.Second * 15)
		}

		// Set processingInProgress to false when processing is complete
		processingInProgress = false
	}()
	fmt.Println("Beginning processing...")
	w.WriteHeader(http.StatusAccepted)
}

func Stats(w http.ResponseWriter, r *http.Request) {
	if processingFinished {
		http.Error(w, "Processing has finished", http.StatusPreconditionFailed)
		return
	}

	if !processingInProgress {
		http.Error(w, "Processing has not started", http.StatusPreconditionFailed)
		return
	}

	processedItems, err := database.GetProcessedItemCount()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calling stats endpoint: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := database.ProcessedItems{Count: processedItems}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func ExtractMiddleName(name string) string {
	// Split the name into words
	words := strings.Fields(name)

	// Check if there are at least three words (first name, middle name, last name)
	if len(words) > 2 {
		// Extract the middle name by joining the words between the first and last names
		middleNameWords := words[1 : len(words)-1]
		middleName := strings.Join(middleNameWords, " ")
		return middleName
	}

	// No middle name present
	return ""
}
