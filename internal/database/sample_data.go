package database

import (
	"data-processing/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var indexModels = []mongo.IndexModel{
	{
		Keys: bson.D{{Key: "_id", Value: 1}},
	},
	{
		Keys: bson.D{{Key: "name", Value: 1}},
	},
	{
		Keys: bson.D{{Key: "age", Value: 1}},
	},
	{
		Keys: bson.D{{Key: "ukBased", Value: 1}},
	},
	{
		Keys: bson.D{{Key: "middleNames", Value: 1}},
	},
	{
		Keys: bson.D{{Key: "processed", Value: 1}},
	},
}

var sampleData = []interface{}{
	models.Person{
		ID:      "1",
		Name:    "Amelia Abigail Alice Adams",
		Age:     10,
		UKBased: true,
	},
	models.Person{
		ID:      "2",
		Name:    "Amelia Abigail Alice Adams",
		Age:     10,
		UKBased: true,
	},
	models.Person{
		ID:      "3",
		Name:    "Clara Catherine Christine Carter",
		Age:     30,
		UKBased: false,
	},
	models.Person{
		ID:      "4",
		Name:    "Daniel David Dominic Douglas",
		Age:     40,
		UKBased: true,
	},
	models.Person{
		ID:      "5",
		Name:    "Emily Elise Eleanor Edwards",
		Age:     50,
		UKBased: false,
	},
	models.Person{
		ID:      "6",
		Name:    "Franklin Felix Frederick Foster",
		Age:     60,
		UKBased: false,
	},
	models.Person{
		ID:      "7",
		Name:    "Grace Gabrielle Genevieve Gray",
		Age:     70,
		UKBased: true,
	},
	models.Person{
		ID:      "8",
		Name:    "Henry Harrison Hayes Hunter",
		Age:     80,
		UKBased: false,
	},
	models.Person{
		ID:      "9",
		Name:    "Isabelle Ivy Irene Ingram",
		Age:     90,
		UKBased: true,
	},
	models.Person{
		ID:      "10",
		Name:    "James Jasper Julian Johnson",
		Age:     100,
		UKBased: false,
	},
}
