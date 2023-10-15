package db

import (
	"context"
	"os"

	"github.com/johnson7543/ims/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const materialColl = "materials"

type MaterialStore interface {
	InsertMaterial(context.Context, *types.Material) (*types.Material, error)
	GetMaterials(context.Context, bson.M) ([]*types.Material, error)
	DeleteMaterial(ctx context.Context, id primitive.ObjectID) (int64, error)
	GetMaterialColors(ctx context.Context) ([]string, error)
	GetMaterialSizes(ctx context.Context) ([]string, error)
}

type MongoMaterialStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoMaterialStore(client *mongo.Client) *MongoMaterialStore {
	dbname := os.Getenv(MongoDBNameEnvName)
	return &MongoMaterialStore{
		client: client,
		coll:   client.Database(dbname).Collection(materialColl),
	}
}

func (s *MongoMaterialStore) GetMaterials(ctx context.Context, filter bson.M) ([]*types.Material, error) {
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

	var materials []*types.Material
	if err := resp.All(ctx, &materials); err != nil {
		return nil, err
	}

	return materials, nil
}

func (s *MongoMaterialStore) InsertMaterial(ctx context.Context, material *types.Material) (*types.Material, error) {
	resp, err := s.coll.InsertOne(ctx, material)
	if err != nil {
		return nil, err
	}
	material.ID = resp.InsertedID.(primitive.ObjectID)

	return material, nil
}

func (s *MongoMaterialStore) DeleteMaterial(ctx context.Context, id primitive.ObjectID) (int64, error) {
	deleteResult, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})

	return deleteResult.DeletedCount, err
}

func (s *MongoMaterialStore) GetMaterialColors(ctx context.Context) ([]string, error) {
	colors, err := s.coll.Distinct(ctx, "color", bson.M{})
	if err != nil {
		return nil, err
	}

	var colorStrings []string
	for _, c := range colors {
		if color, ok := c.(string); ok {
			colorStrings = append(colorStrings, color)
		}
	}

	return colorStrings, nil
}

func (s *MongoMaterialStore) GetMaterialSizes(ctx context.Context) ([]string, error) {
	sizes, err := s.coll.Distinct(ctx, "size", bson.M{})
	if err != nil {
		return nil, err
	}

	var sizeStrings []string
	for _, size := range sizes {
		if sizeStr, ok := size.(string); ok {
			sizeStrings = append(sizeStrings, sizeStr)
		}
	}

	return sizeStrings, nil
}
