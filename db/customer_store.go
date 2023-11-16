package db

import (
	"context"
	"os"

	"github.com/johnson7543/ims/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const customerColl = "customers"

type CustomerStore interface {
	GetCustomers(context.Context, bson.M) ([]*types.Customer, error)
	InsertCustomer(context.Context, *types.Customer) (*types.Customer, error)
	UpdateCustomer(ctx context.Context, id primitive.ObjectID, updatedCustomer *types.Customer) (int64, error)
	DeleteCustomer(ctx context.Context, id primitive.ObjectID) (int64, error)
}

type MongoCustomerStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoCustomerStore(client *mongo.Client) *MongoCustomerStore {
	dbname := os.Getenv(MongoDBNameEnvName)
	return &MongoCustomerStore{
		client: client,
		coll:   client.Database(dbname).Collection(customerColl),
	}
}

func (s *MongoCustomerStore) GetCustomers(ctx context.Context, filter bson.M) ([]*types.Customer, error) {
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

	var customers []*types.Customer
	if err := resp.All(ctx, &customers); err != nil {
		return nil, err
	}

	return customers, nil
}

func (s *MongoCustomerStore) InsertCustomer(ctx context.Context, customer *types.Customer) (*types.Customer, error) {
	resp, err := s.coll.InsertOne(ctx, customer)
	if err != nil {
		return nil, err
	}
	customer.ID = resp.InsertedID.(primitive.ObjectID)

	return customer, nil
}

func (s *MongoCustomerStore) UpdateCustomer(ctx context.Context, id primitive.ObjectID, updatedCustomer *types.Customer) (int64, error) {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"company":     updatedCustomer.Company,
			"name":        updatedCustomer.Name,
			"phone":       updatedCustomer.Phone,
			"address":     updatedCustomer.Address,
			"taxIdNumber": updatedCustomer.TaxIdNumber,
		},
	}

	updateResult, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return updateResult.ModifiedCount, nil
}

func (s *MongoCustomerStore) DeleteCustomer(ctx context.Context, id primitive.ObjectID) (int64, error) {
	deleteResult, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	return deleteResult.DeletedCount, err
}
