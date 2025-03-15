package db

import (
	indexes "github.com/shubhamoys/todo-api/db/db_indexes"
	"go.mongodb.org/mongo-driver/mongo"
)

// EnsureIndexes calls all index functions for different collections
func EnsureIndexes(db *mongo.Database) {
	indexes.EnsureUserIndexes(db)
	// Add more index functions as needed
}