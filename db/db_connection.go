package db

import (
	"context"
	"time"

	"github.com/shubhamoys/todo-api/config"
	"github.com/shubhamoys/todo-api/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect(mongoURI string, dbName string) {
	// Set client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		utils.Logger.Error("Failed to connect to MongoDB: ", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		utils.Logger.Error("Failed to ping MongoDB: ", err)
	}

	utils.Logger.Info("Connected to MongoDB!")
	Client = client

	DB := client.Database(dbName)

	EnsureIndexes(DB)
}

func GetCollection(collectionName string) *mongo.Collection {
	dbName := config.AppConfig.DBName
	return Client.Database(dbName).Collection(collectionName)
}

func Disconnect() {
	if Client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := Client.Disconnect(ctx)
		if err != nil {
			utils.Logger.Error("Failed to disconnect from MongoDB: ", err)
		}

		utils.Logger.Info("Successfully Disconnected from MongoDB!")
	}
}
