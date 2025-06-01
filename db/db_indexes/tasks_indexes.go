package indexes

import (
	"context"

	"github.com/shubhamoys/todo-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// EnsureTaskIndexes creates required indexes on the tasks collection
func EnsureTaskIndexes(db *mongo.Database) {
	tasksCollection := db.Collection("tasks")

	_, err := tasksCollection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "name", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "description", Value: 1}},
		},
	})
	if err != nil {
		utils.Logger.Error("Error creating task indexes:", err)
	}

	utils.Logger.Info("task indexes created successfully")
}
