package config

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoConfig struct {
	URI      string        `json:"uri"`
	Database string        `json:"database"`
	Timeout  time.Duration `json:"timeout"`
}

var mongoClient *mongo.Client

func InitMongoDB(config MongoConfig) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.URI).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(opts)
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to MongoDB: %v", err))
	}

	// Verify connection with ping
	var result bson.M
	if err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		panic(fmt.Sprintf("Failed to ping MongoDB: %v", err))
	}

	// Store client globally for later cleanup
	mongoClient = client

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return client
}

func GetMongoClient() *mongo.Client {
	if mongoClient == nil {
		panic("MongoDB client is not initialized. Call InitMongoDB first.")
	}
	return mongoClient
}

// CloseMongoDB closes the MongoDB connection gracefully
func CloseMongoDB() error {
	if mongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := mongoClient.Disconnect(ctx); err != nil {
			return fmt.Errorf("failed to disconnect MongoDB: %w", err)
		}
		mongoClient = nil
		fmt.Println("MongoDB connection closed successfully")
	}
	return nil
}
