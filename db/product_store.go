package db

import (
	"context"
	"os"

	"github.com/johnson7543/ims/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const productColl = "products"

type ProductStore interface {
	InsertProduct(context.Context, *types.Product) (*types.Product, error)
	GetProducts(context.Context, bson.M) ([]*types.Product, error)
	DeleteProduct(ctx context.Context, id primitive.ObjectID) error
	GetProductColors(ctx context.Context) ([]string, error)
	GetProductSizes(ctx context.Context) ([]string, error)
}

type MongoProductStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoProductStore(client *mongo.Client) *MongoProductStore {
	dbName := os.Getenv(MongoDBNameEnvName)
	return &MongoProductStore{
		client: client,
		coll:   client.Database(dbName).Collection(productColl),
	}
}

func (s *MongoProductStore) GetProducts(ctx context.Context, filter bson.M) ([]*types.Product, error) {
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

	var products []*types.Product
	if err := resp.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (s *MongoProductStore) InsertProduct(ctx context.Context, product *types.Product) (*types.Product, error) {
	resp, err := s.coll.InsertOne(ctx, product)
	if err != nil {
		return nil, err
	}
	product.ID = resp.InsertedID.(primitive.ObjectID)

	return product, nil
}

func (s *MongoProductStore) DeleteProduct(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (s *MongoProductStore) GetProductColors(ctx context.Context) ([]string, error) {
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

func (s *MongoProductStore) GetProductSizes(ctx context.Context) ([]string, error) {
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