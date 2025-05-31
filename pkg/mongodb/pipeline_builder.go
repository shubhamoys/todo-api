package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// LookupConfig defines the configuration for a $lookup stage
// Used to specify how collections should be joined
type LookupConfig struct {
	From         string // Target collection to join with
	LocalField   string // Field from the current collection
	ForeignField string // Field from the target collection to match against
	As           string // Name of the output array field
}

// PipelineBuilder implements the Builder pattern for MongoDB aggregation pipelines
type PipelineBuilder struct {
	pipeline mongo.Pipeline
}

// NewPipelineBuilder creates a new instance of PipelineBuilder
func NewPipelineBuilder() *PipelineBuilder {
	return &PipelineBuilder{
		pipeline: mongo.Pipeline{},
	}
}

// AddMatch adds a $match stage to filter documents
// $match: Filters the documents to pass only those that match the specified condition(s)
func (pb *PipelineBuilder) AddMatch(query bson.M) *PipelineBuilder {
	pb.pipeline = append(pb.pipeline, bson.D{{Key: "$match", Value: query}})
	return pb
}

// AddSort adds a $sort stage to order documents
// $sort: Sorts all input documents and returns them in the specified order
// Example: bson.D{{"created_at", -1}} for descending order
func (pb *PipelineBuilder) AddSort(sort bson.D) *PipelineBuilder {
	pb.pipeline = append(pb.pipeline, bson.D{{Key: "$sort", Value: sort}})
	return pb
}

// AddPagination adds $skip and $limit stages for pagination
// $skip: Skips the first n documents where n is the specified skip number
// $limit: Limits the number of documents to the specified number
func (pb *PipelineBuilder) AddPagination(skip, limit int64) *PipelineBuilder {
	if skip > 0 {
		pb.pipeline = append(pb.pipeline, bson.D{{Key: "$skip", Value: skip}})
	}
	if limit > 0 {
		pb.pipeline = append(pb.pipeline, bson.D{{Key: "$limit", Value: limit}})
	}
	return pb
}

// AddLookup adds a $lookup stage to join with another collection
// $lookup: Performs a left outer join to another collection
// Automatically adds $unwind to flatten the resulting array
func (pb *PipelineBuilder) AddLookup(config LookupConfig) *PipelineBuilder {
	// Validate config
	if config.From == "" || config.LocalField == "" ||
		config.ForeignField == "" || config.As == "" {
		return pb // Skip invalid lookup config
	}

	// Add lookup stage
	pb.pipeline = append(pb.pipeline, bson.D{{Key: "$lookup", Value: bson.M{
		"from":         config.From,
		"localField":   config.LocalField,
		"foreignField": config.ForeignField,
		"as":           config.As,
	}}})

	// Add unwind stage
	// $unwind: Deconstructs an array field to output one document for each element
	pb.pipeline = append(pb.pipeline, bson.D{{Key: "$unwind", Value: bson.M{
		"path":                       "$" + config.As,
		"preserveNullAndEmptyArrays": true, // Keep documents that don't match
	}}})

	return pb
}

// AddProjection adds a $project stage to shape the output documents
// $project: Reshapes documents by specifying which fields to include/exclude
// Example: bson.M{"_id": 0, "name": 1} to exclude _id and include name
func (pb *PipelineBuilder) AddProjection(projection bson.M) *PipelineBuilder {
	if len(projection) > 0 {
		pb.pipeline = append(pb.pipeline, bson.D{{Key: "$project", Value: projection}})
	}
	return pb
}

// Build returns the final MongoDB aggregation pipeline
func (pb *PipelineBuilder) Build() mongo.Pipeline {
	return pb.pipeline
}
