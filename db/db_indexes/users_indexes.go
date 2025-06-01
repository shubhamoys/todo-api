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

	indexes := []mongo.IndexModel{
		{
			// Unique email index
			Keys: bson.D{{Key: "email.value", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetSparse(true),
		},
		{
			// Index for name search
			Keys: bson.D{{Key: "name", Value: 1}},
		},
		{
			// Index for email search
			Keys: bson.D{{Key: "email.value", Value: 1}},
		},
	}

	_, err := usersCollection.Indexes().CreateMany(context.Background(), indexes)
	if err != nil {
		utils.Logger.Error("Error creating user indexes:", err)
	}

	utils.Logger.Info("User indexes created successfully")
}
