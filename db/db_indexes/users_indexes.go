package indexes

import (
	"context"

	"github.com/shubhamoys/todo-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// EnsureUserIndexes creates required indexes on the users collection
func EnsureUserIndexes(db *mongo.Database) {
	usersCollection := db.Collection("users")

	uniqueEmailIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "email.value", Value: 1}}, // Unique index on email.value
		Options: options.Index().
			SetUnique(true).
			SetSparse(true),
	}

	// Full-text index on name and email.value
	textIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "name", Value: "text"},
			{Key: "email.value", Value: "text"},
		},
	}

	_, err := usersCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{uniqueEmailIndex, textIndex})
	if err != nil {
		utils.Logger.Error("Error creating user indexes:", err)
	}

	utils.Logger.Info("User indexes created successfully")
}
