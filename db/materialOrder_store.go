package db

import (
	"context"
	"os"

	"github.com/johnson7543/ims/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const materialOrderColl = "materialOrders"

type MaterialOrderStore interface {
	GetMaterialOrders(ctx context.Context, filter bson.M) ([]*types.MaterialOrder, error)
	InsertMaterialOrder(ctx context.Context, order *types.MaterialOrder) (*types.MaterialOrder, error)
	UpdateMaterialOrder(ctx context.Context, orderID primitive.ObjectID, updatedOrder *types.MaterialOrder) (int64, error)
	DeleteMaterialOrder(ctx context.Context, id primitive.ObjectID) (int64, error)
}

type MongoMaterialOrderStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoMaterialOrderStore(client *mongo.Client) *MongoMaterialOrderStore {
	dbname := os.Getenv(MongoDBNameEnvName)
	return &MongoMaterialOrderStore{
		client: client,
		coll:   client.Database(dbname).Collection(materialOrderColl),
	}
}

func (s *MongoMaterialOrderStore) GetMaterialOrders(ctx context.Context, filter bson.M) ([]*types.MaterialOrder, error) {
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

	var materialOrders []*types.MaterialOrder
	if err := resp.All(ctx, &materialOrders); err != nil {
		return nil, err
	}

	return materialOrders, nil
}

func (s *MongoMaterialOrderStore) InsertMaterialOrder(ctx context.Context, order *types.MaterialOrder) (*types.MaterialOrder, error) {
	resp, err := s.coll.InsertOne(ctx, order)
	if err != nil {
		return nil, err
	}
	order.ID = resp.InsertedID.(primitive.ObjectID)

	return order, nil
}

func (s *MongoMaterialOrderStore) UpdateMaterialOrder(ctx context.Context, orderID primitive.ObjectID, updatedOrder *types.MaterialOrder) (int64, error) {
	filter := bson.M{"_id": orderID}

	update := bson.M{
		"$set": bson.M{
			"sellerId":     updatedOrder.SellerID,
			"orderDate":    updatedOrder.OrderDate,
			"deliveryDate": updatedOrder.DeliveryDate,
			"paymentDate":  updatedOrder.PaymentDate,
			"totalAmount":  updatedOrder.TotalAmount,
			"status":       updatedOrder.Status,
		},
	}

	updateResult, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return updateResult.ModifiedCount, nil
}

func (s *MongoMaterialOrderStore) DeleteMaterialOrder(ctx context.Context, id primitive.ObjectID) (int64, error) {
	deleteResult, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	return deleteResult.DeletedCount, err
}
