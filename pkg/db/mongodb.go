package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB holds the MongoDB client and database
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewMongoDB creates a new MongoDB connection
func NewMongoDB(uri, dbName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	log.Println("Connected to MongoDB successfully")

	db := client.Database(dbName)

	// Create indexes
	if err := createIndexes(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to create indexes: %w", err)
	}

	return &MongoDB{
		Client:   client,
		Database: db,
	}, nil
}

// createIndexes creates necessary indexes for the collections
func createIndexes(ctx context.Context, db *mongo.Database) error {
	// Users indexes
	usersIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	if _, err := db.Collection("users").Indexes().CreateMany(ctx, usersIndexes); err != nil {
		return err
	}

	// Monthly records indexes
	monthlyIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "year", Value: 1},
				{Key: "month", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
	if _, err := db.Collection("monthly_records").Indexes().CreateMany(ctx, monthlyIndexes); err != nil {
		return err
	}

	// Yearly summaries indexes
	yearlyIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "year", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	}
	if _, err := db.Collection("yearly_summaries").Indexes().CreateMany(ctx, yearlyIndexes); err != nil {
		return err
	}

	// Default commitments indexes
	defaultCommitmentsIndexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}
	if _, err := db.Collection("default_commitments").Indexes().CreateMany(ctx, defaultCommitmentsIndexes); err != nil {
		return err
	}

	log.Println("Database indexes created successfully")
	return nil
}

// Close closes the MongoDB connection
func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return m.Client.Disconnect(ctx)
}

// Collection returns a MongoDB collection
func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.Database.Collection(name)
}
