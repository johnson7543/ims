package db

import (
	"context"
	"os"

	"github.com/johnson7543/ims/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const sellerColl = "sellers"

type SellerStore interface {
	GetSellers(context.Context, bson.M) ([]*types.Seller, error)
	InsertSeller(context.Context, *types.Seller) (*types.Seller, error)
	UpdateSeller(context.Context, primitive.ObjectID, *types.Seller) (int64, error)
	DeleteSeller(context.Context, primitive.ObjectID) (int64, error)
}

type MongoSellerStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoSellerStore(client *mongo.Client) *MongoSellerStore {
	dbname := os.Getenv(MongoDBNameEnvName)
	return &MongoSellerStore{
		client: client,
		coll:   client.Database(dbname).Collection(sellerColl),
	}
}

func (s *MongoSellerStore) GetSellers(ctx context.Context, filter bson.M) ([]*types.Seller, error) {
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

	var sellers []*types.Seller
	if err := resp.All(ctx, &sellers); err != nil {
		return nil, err
	}

	return sellers, nil
}

func (s *MongoSellerStore) InsertSeller(ctx context.Context, seller *types.Seller) (*types.Seller, error) {
	resp, err := s.coll.InsertOne(ctx, seller)
	if err != nil {
		return nil, err
	}
	seller.ID = resp.InsertedID.(primitive.ObjectID)

	return seller, nil
}

func (s *MongoSellerStore) UpdateSeller(ctx context.Context, id primitive.ObjectID, updatedSeller *types.Seller) (int64, error) {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"company":     updatedSeller.Company,
			"name":        updatedSeller.Name,
			"phone":       updatedSeller.Phone,
			"address":     updatedSeller.Address,
			"taxIdNumber": updatedSeller.TaxIdNumber,
		},
	}

	updateResult, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return updateResult.ModifiedCount, nil
}

func (s *MongoSellerStore) DeleteSeller(ctx context.Context, id primitive.ObjectID) (int64, error) {
	deleteResult, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	return deleteResult.DeletedCount, err
}
