package db

import (
	"context"
	"fmt"
	"time"

	models "movie_booking_system/model"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
)

func ConnectDB(ctx context.Context, config models.Config) error {
	clientOptions := options.Client().ApplyURI(config.MongoURI).
		SetConnectTimeout(10 * time.Second).
		SetServerSelectionTimeout(30 * time.Second)

	connectCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var err error
	client, err = mongo.Connect(connectCtx, clientOptions)
	if err != nil {
		return fmt.Errorf("connect mongo client: %w", err)
	}

	pingCtx, pingCancel := context.WithTimeout(ctx, 5*time.Second)
	defer pingCancel()

	err = client.Ping(pingCtx, nil)
	if err != nil {
		return fmt.Errorf("ping mongo database: %w", err)
	}

	database = client.Database(config.MongoDatabase)
	collection = database.Collection(config.MongoCollection)

	return nil
}

func GetCollection() *mongo.Collection {
	return collection
}

func Ping(ctx context.Context) error {
	if client == nil {
		return fmt.Errorf("database client is not initialized")
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		return fmt.Errorf("ping mongo database: %w", err)
	}

	return nil
}

func Close(ctx context.Context) error {
	if client == nil {
		return nil
	}

	if err := client.Disconnect(ctx); err != nil {
		return fmt.Errorf("disconnect mongo client: %w", err)
	}

	return nil
}
