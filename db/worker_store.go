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
	InsertWorker(context.Context, *types.Worker) (*types.Worker, error)
	GetWorkers(context.Context, bson.M) ([]*types.Worker, error)
	DeleteWorker(ctx context.Context, id primitive.ObjectID) error
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

func (s *MongoWorkerStore) DeleteWorker(ctx context.Context, id primitive.ObjectID) error {
	_, err := s.coll.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
