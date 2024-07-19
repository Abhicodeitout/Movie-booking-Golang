package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
)

func ConnectDB() error {
	// Replace with your MongoDB Atlas connection string
	uri := "mongodb+srv://<username>:<password>@cluster0.shtcmvu.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

	// Set client options
	clientOptions := options.Client().ApplyURI(uri).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(30 * time.Second)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Ping the MongoDB server
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Set database and collection
	database = client.Database("movie_booking")
	collection = database.Collection("movies")

	return nil
}

func GetCollection() *mongo.Collection {
	return collection
}
