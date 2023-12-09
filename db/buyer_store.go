package db

import (
	"context"
	"os"

	"github.com/johnson7543/ims/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const buyerColl = "buyers"

type BuyerStore interface {
	GetBuyers(context.Context, bson.M) ([]*types.Buyer, error)
	InsertBuyer(context.Context, *types.Buyer) (*types.Buyer, error)
	UpdateBuyer(ctx context.Context, id primitive.ObjectID, updatedBuyer *types.Buyer) (int64, error)
	DeleteBuyer(ctx context.Context, id primitive.ObjectID) (int64, error)
}

type MongoBuyerStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBuyerStore(client *mongo.Client) *MongoBuyerStore {
	dbname := os.Getenv(MongoDBNameEnvName)
	return &MongoBuyerStore{
		client: client,
		coll:   client.Database(dbname).Collection(buyerColl),
	}
}

func (s *MongoBuyerStore) GetBuyers(ctx context.Context, filter bson.M) ([]*types.Buyer, error) {
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

	var buyers []*types.Buyer
	if err := resp.All(ctx, &buyers); err != nil {
		return nil, err
	}

	return buyers, nil
}

func (s *MongoBuyerStore) InsertBuyer(ctx context.Context, buyer *types.Buyer) (*types.Buyer, error) {
	resp, err := s.coll.InsertOne(ctx, buyer)
	if err != nil {
		return nil, err
	}
	buyer.ID = resp.InsertedID.(primitive.ObjectID)

	return buyer, nil
}

func (s *MongoBuyerStore) UpdateBuyer(ctx context.Context, id primitive.ObjectID, updatedBuyer *types.Buyer) (int64, error) {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"company":     updatedBuyer.Company,
			"name":        updatedBuyer.Name,
			"phone":       updatedBuyer.Phone,
			"address":     updatedBuyer.Address,
			"taxIdNumber": updatedBuyer.TaxIdNumber,
		},
	}

	updateResult, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return updateResult.ModifiedCount, nil
}

func (s *MongoBuyerStore) DeleteBuyer(ctx context.Context, id primitive.ObjectID) (int64, error) {
	deleteResult, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	return deleteResult.DeletedCount, err
}
