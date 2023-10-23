package db

import (
	"context"
	"os"

	"github.com/johnson7543/ims/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const workerColl = "workers"

type WorkerStore interface {
	GetWorkers(context.Context, bson.M) ([]*types.Worker, error)
	InsertWorker(context.Context, *types.Worker) (*types.Worker, error)
	UpdateWorker(ctx context.Context, id primitive.ObjectID, updatedWorker *types.Worker) (int64, error)
	DeleteWorker(ctx context.Context, id primitive.ObjectID) (int64, error)
}

type MongoWorkerStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoWorkerStore(client *mongo.Client) *MongoWorkerStore {
	dbname := os.Getenv(MongoDBNameEnvName)
	return &MongoWorkerStore{
		client: client,
		coll:   client.Database(dbname).Collection(workerColl),
	}
}

func (s *MongoWorkerStore) GetWorkers(ctx context.Context, filter bson.M) ([]*types.Worker, error) {
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

	var workers []*types.Worker
	if err := resp.All(ctx, &workers); err != nil {
		return nil, err
	}

	return workers, nil
}

func (s *MongoWorkerStore) InsertWorker(ctx context.Context, worker *types.Worker) (*types.Worker, error) {
	resp, err := s.coll.InsertOne(ctx, worker)
	if err != nil {
		return nil, err
	}
	worker.ID = resp.InsertedID.(primitive.ObjectID)

	return worker, nil
}

func (s *MongoWorkerStore) UpdateWorker(ctx context.Context, id primitive.ObjectID, updatedWorker *types.Worker) (int64, error) {
	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"company":     updatedWorker.Company,
			"name":        updatedWorker.Name,
			"phone":       updatedWorker.Phone,
			"address":     updatedWorker.Address,
			"taxIdNumber": updatedWorker.TaxIdNumber,
		},
	}

	updateResult, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return updateResult.ModifiedCount, nil
}

func (s *MongoWorkerStore) DeleteWorker(ctx context.Context, id primitive.ObjectID) (int64, error) {
	deleteResult, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	return deleteResult.DeletedCount, err
}
