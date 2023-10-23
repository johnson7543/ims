package db

import (
	"context"
	"os"

	"github.com/johnson7543/ims/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const orderColl = "orders"

type OrderStore interface {
	GetOrders(ctx context.Context, filter bson.M) ([]*types.Order, error)
	InsertOrder(ctx context.Context, order *types.Order) (*types.Order, error)
	UpdateOrder(ctx context.Context, orderID primitive.ObjectID, updatedOrder *types.Order) (int64, error)
	DeleteOrder(ctx context.Context, id primitive.ObjectID) (int64, error)
}

type MongoOrderStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoOrderStore(client *mongo.Client) *MongoOrderStore {
	dbname := os.Getenv(MongoDBNameEnvName)
	return &MongoOrderStore{
		client: client,
		coll:   client.Database(dbname).Collection(orderColl),
	}
}

func (s *MongoOrderStore) GetOrders(ctx context.Context, filter bson.M) ([]*types.Order, error) {
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

	var orders []*types.Order
	if err := resp.All(ctx, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *MongoOrderStore) InsertOrder(ctx context.Context, order *types.Order) (*types.Order, error) {
	resp, err := s.coll.InsertOne(ctx, order)
	if err != nil {
		return nil, err
	}
	order.ID = resp.InsertedID.(primitive.ObjectID)

	return order, nil
}

func (s *MongoOrderStore) UpdateOrder(ctx context.Context, orderID primitive.ObjectID, updatedOrder *types.Order) (int64, error) {
	filter := bson.M{"_id": orderID}
	update := bson.M{
		"$set": bson.M{
			"totalAmount":     updatedOrder.TotalAmount,
			"status":          updatedOrder.Status,
			"shippingAddress": updatedOrder.ShippingAddress,
			"orderItems":      updatedOrder.OrderItems,
		},
	}

	updateResult, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return updateResult.ModifiedCount, nil
}

func (s *MongoOrderStore) DeleteOrder(ctx context.Context, id primitive.ObjectID) (int64, error) {
	deleteResult, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	return deleteResult.DeletedCount, err
}
