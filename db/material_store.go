package db

import (
	"context"
	"errors"
	"os"

	"github.com/johnson7543/ims/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const materialColl = "materials"

type MaterialStore interface {
	GetMaterials(context.Context, bson.M) ([]*types.Material, error)
	GetMaterial(context.Context, primitive.ObjectID) (*types.Material, error)
	InsertMaterial(context.Context, *types.Material) (*types.Material, error)
	UpdateMaterial(context.Context, primitive.ObjectID, *types.Material) (int64, error)
	DeleteMaterial(context.Context, primitive.ObjectID) (int64, error)
	GetMaterialColors(context.Context, string) ([]string, error)
	GetMaterialTypes(context.Context) ([]string, error)
	GetMaterialSizes(context.Context, string) ([]string, error)
	DecreaseMaterialQuantity(ctx context.Context, materialID primitive.ObjectID, quantity int) (int64, error)
	IncreaseMaterialQuantity(ctx context.Context, materialID primitive.ObjectID, quantity int) (int64, error)
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

func (s *MongoMaterialStore) GetMaterial(ctx context.Context, materialID primitive.ObjectID) (*types.Material, error) {
	filter := bson.M{"_id": materialID}
	var material types.Material

	err := s.coll.FindOne(ctx, filter).Decode(&material)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("material not found")
		}
		return nil, err
	}

	return &material, nil
}

func (s *MongoMaterialStore) InsertMaterial(ctx context.Context, material *types.Material) (*types.Material, error) {
	resp, err := s.coll.InsertOne(ctx, material)
	if err != nil {
		return nil, err
	}
	material.ID = resp.InsertedID.(primitive.ObjectID)

	return material, nil
}

func (s *MongoMaterialStore) UpdateMaterial(ctx context.Context, materialID primitive.ObjectID, updates *types.Material) (int64, error) {
	filter := bson.M{"_id": materialID}
	update := bson.M{
		"$set": bson.M{
			"name":          updates.Name,
			"color":         updates.Color,
			"type":          updates.Type,
			"size":          updates.Size,
			"quantity":      updates.Quantity,
			"remarks":       updates.Remarks,
			"price_history": updates.PriceHistory,
		},
	}

	updateResult, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return updateResult.ModifiedCount, nil
}

func (s *MongoMaterialStore) DeleteMaterial(ctx context.Context, id primitive.ObjectID) (int64, error) {
	deleteResult, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})

	return deleteResult.DeletedCount, err
}

func (s *MongoMaterialStore) GetMaterialColors(ctx context.Context, materialType string) ([]string, error) {
	filter := bson.M{}
	if materialType != "" {
		filter["type"] = materialType
	}

	colors, err := s.coll.Distinct(ctx, "color", filter)
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

func (s *MongoMaterialStore) GetMaterialTypes(ctx context.Context) ([]string, error) {
	types, err := s.coll.Distinct(ctx, "type", bson.M{})
	if err != nil {
		return nil, err
	}

	var typeStrings []string
	for _, t := range types {
		if typeStr, ok := t.(string); ok {
			typeStrings = append(typeStrings, typeStr)
		}
	}

	return typeStrings, nil
}

func (s *MongoMaterialStore) GetMaterialSizes(ctx context.Context, materialType string) ([]string, error) {
	filter := bson.M{}
	if materialType != "" {
		filter["type"] = materialType
	}

	sizes, err := s.coll.Distinct(ctx, "size", filter)
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

func (s *MongoMaterialStore) DecreaseMaterialQuantity(ctx context.Context, materialID primitive.ObjectID, quantity int) (int64, error) {
	filter := bson.M{"_id": materialID, "quantity": bson.M{"$gte": quantity}}
	update := bson.M{
		"$inc": bson.M{"quantity": -quantity},
	}

	updateResult, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return updateResult.ModifiedCount, nil
}

func (s *MongoMaterialStore) IncreaseMaterialQuantity(ctx context.Context, materialID primitive.ObjectID, quantity int) (int64, error) {
	filter := bson.M{"_id": materialID}
	update := bson.M{
		"$inc": bson.M{"quantity": quantity},
	}

	updateResult, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return updateResult.ModifiedCount, nil
}
