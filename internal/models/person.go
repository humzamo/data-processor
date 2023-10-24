package models

// Person is a struct for the documents in the Persons collection
type Person struct {
	ID          string `bson:"_id" json:"id"`
	Name        string `bson:"name" json:"name"`
	Age         int    `bson:"age" json:"age"`
	UKBased     bool   `bson:"ukBased" json:"ukBased"`
	MiddleNames string `bson:"middleNames,omitempty" json:"middleNames,omitempty"`
	Processed   bool   `bson:"processed" json:"processed"`
}
