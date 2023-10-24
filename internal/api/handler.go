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
	status ProcessingStatus
)

const (
	batchSize = 3
	sleepTime = time.Second * 15
)

// StartHandler is the handler for the start endpoint
func StartHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Call to StartHandler")

	if status == ProcessingStatusFinished {
		http.Error(w, "Processing has already finished!", http.StatusAccepted)
		return
	}

	if status == ProcessingStatusProcessing {
		http.Error(w, "Processing is already in progress", http.StatusTooManyRequests)
		return
	}

	// Set processing in progress and start processing
	status = ProcessingStatusProcessing
	fmt.Println("Processing...")

	go func() {
		for {
			if status == ProcessingStatusPaused {
				fmt.Println("Processing has been paused")
				return
			}

			// Get a batch of unprocessed persons from the database
			persons, err := database.GetUnprocessedPersons(batchSize)
			if err != nil {
				http.Error(w, "Error fetching unprocessed persons", http.StatusInternalServerError)
				fmt.Println("Error fetching unprocessed persons. Pausing processing.", err.Error())
				status = ProcessingStatusPaused
				return
			}

			// No more unprocessed persons, stop processing
			if len(persons) == 0 {
				break
			}

			ctx, cancel := context.WithTimeout(context.Background(), database.DefaultTimeout)
			defer cancel()

			collection := database.DB.Collection(database.PersonsCollection)

			// Process the batch of persons
			for _, person := range persons {
				person.MiddleNames = extractMiddleName(person.Name)

				// Update the "processed" flag to mark it as processed
				person.Processed = true

				_, err := collection.ReplaceOne(ctx, bson.M{"_id": person.ID}, person)
				if err != nil {
					log.Printf("Error updating document ID %s for person: %v", person.ID, err)
				}
				log.Printf("Updated document ID %s", person.ID)
			}

			// Artifical sleep added for testing purposes
			fmt.Println("Sleeping...")
			time.Sleep(sleepTime)
		}

		// Set status to finished when processing is complete
		fmt.Println("Processing has finshed!")
		status = ProcessingStatusFinished
	}()

	w.WriteHeader(http.StatusAccepted)
}

// StatsHandler is the handler for the stats endpoint
func StatsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Call to StatsHandler")

	if status == ProcessingStatusNotStarted {
		http.Error(w, "Processing has not started", http.StatusPreconditionFailed)
		return
	}

	processedItems, err := database.GetProcessedItemCount()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calling stats endpoint: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response := database.ProcessedItems{Count: processedItems, Status: string(status)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Pause is the handler for the pause endpoint
func PauseHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Call to PauseHandler")

	if status == ProcessingStatusFinished {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Processing has already finished"))
		return
	}

	if status == ProcessingStatusNotStarted {
		http.Error(w, "Processing has not started", http.StatusPreconditionFailed)
		return
	}

	status = ProcessingStatusPaused
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Processing has been paused"))
}

func extractMiddleName(name string) string {
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
