package db

import (
	"context"
	"os"

	"github.com/johnson7543/ims/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const processingItemColl = "processing_items"

type ProcessingItemStore interface {
	InsertProcessingItem(context.Context, *types.ProcessingItem) (*types.ProcessingItem, error)
	GetProcessingItems(context.Context, bson.M) ([]*types.ProcessingItem, error)
	DeleteProcessingItem(ctx context.Context, id primitive.ObjectID) error
}

type MongoProcessingItemStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoProcessingItemStore(client *mongo.Client) *MongoProcessingItemStore {
	dbName := os.Getenv(MongoDBNameEnvName)
	return &MongoProcessingItemStore{
		client: client,
		coll:   client.Database(dbName).Collection(processingItemColl),
	}
}

func (s *MongoProcessingItemStore) GetProcessingItems(ctx context.Context, filter bson.M) ([]*types.ProcessingItem, error) {
	for key, value := range filter {
		switch v := value.(type) {
		case string:
			filter[key] = bson.M{"$regex": primitive.Regex{Pattern: v, Options: "i"}}
		case int, int32, int64, float32, float64:
			filter[key] = bson.M{"$eq": v}
		}
	}

	resp, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var processingItems []*types.ProcessingItem
	if err := resp.All(ctx, &processingItems); err != nil {
		return nil, err
	}

	return processingItems, nil
}

func (s *MongoProcessingItemStore) InsertProcessingItem(ctx context.Context, processingItem *types.ProcessingItem) (*types.ProcessingItem, error) {
	resp, err := s.coll.InsertOne(ctx, processingItem)
	if err != nil {
		return nil, err
	}
	processingItem.ID = resp.InsertedID.(primitive.ObjectID)

	return processingItem, nil
}

func (s *MongoProcessingItemStore) DeleteProcessingItem(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
