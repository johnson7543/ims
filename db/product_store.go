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
	GetProducts(context.Context, bson.M) ([]*types.Product, error)
	InsertProduct(context.Context, *types.Product) (*types.Product, error)
	UpdateProduct(ctx context.Context, productID primitive.ObjectID, updatedProduct *types.Product) (int64, error)
	DeleteProduct(ctx context.Context, id primitive.ObjectID) (int64, error)
	GetProductColors(context.Context, string) ([]string, error)
	GetProductTypes(context.Context) ([]string, error)
	GetProductSizes(context.Context, string) ([]string, error)
	CheckExistedSKU(ctx context.Context, sku string) (bool, error)
	CheckDuplicateSKU(ctx context.Context, sku string, productID primitive.ObjectID) (bool, error)
	DecreaseProductQuantity(ctx context.Context, productID primitive.ObjectID, quantity int) (int64, error)
	IncreaseProductQuantity(ctx context.Context, productID primitive.ObjectID, quantity int) (int64, error)
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

func (s *MongoProductStore) UpdateProduct(ctx context.Context, productID primitive.ObjectID, updatedProduct *types.Product) (int64, error) {
	filter := bson.M{"_id": productID}
	update := bson.M{
		"$set": bson.M{
			"sku":      updatedProduct.SKU,
			"name":     updatedProduct.Name,
			"material": updatedProduct.Material,
			"color":    updatedProduct.Color,
			"type":     updatedProduct.Type,
			"size":     updatedProduct.Size,
			"quantity": updatedProduct.Quantity,
			"price":    updatedProduct.Price,
			"date":     updatedProduct.Date,
			"remark":   updatedProduct.Remark,
		},
	}

	updateResult, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return updateResult.ModifiedCount, nil
}

func (s *MongoProductStore) DeleteProduct(ctx context.Context, id primitive.ObjectID) (int64, error) {
	deleteResult, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	return deleteResult.DeletedCount, err
}

func (s *MongoProductStore) GetProductColors(ctx context.Context, productType string) ([]string, error) {
	filter := bson.M{}
	if productType != "" {
		filter["type"] = productType
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

func (s *MongoProductStore) GetProductTypes(ctx context.Context) ([]string, error) {
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

func (s *MongoProductStore) GetProductSizes(ctx context.Context, productType string) ([]string, error) {
	filter := bson.M{}
	if productType != "" {
		filter["type"] = productType
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

func (s *MongoProductStore) CheckExistedSKU(ctx context.Context, sku string) (bool, error) {
	filter := bson.M{
		"sku": sku,
	}

	count, err := s.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *MongoProductStore) CheckDuplicateSKU(ctx context.Context, sku string, productID primitive.ObjectID) (bool, error) {
	filter := bson.M{
		"_id": bson.M{"$ne": productID}, // Exclude the current product
		"sku": sku,
	}

	count, err := s.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *MongoProductStore) DecreaseProductQuantity(ctx context.Context, productID primitive.ObjectID, quantity int) (int64, error) {
	filter := bson.M{"_id": productID}

	update := bson.M{
		"$inc": bson.M{"quantity": -quantity},
	}

	updateResult, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return updateResult.ModifiedCount, nil
}

func (s *MongoProductStore) IncreaseProductQuantity(ctx context.Context, productID primitive.ObjectID, quantity int) (int64, error) {
	filter := bson.M{"_id": productID}
	update := bson.M{
		"$inc": bson.M{"quantity": quantity},
	}

	updateResult, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return updateResult.ModifiedCount, nil
}
