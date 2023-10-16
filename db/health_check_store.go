package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const adminDB = "admin"

type HealthCheckResponse struct {
	Version string `json:"version"`
	Ok      int    `json:"ok"`
}

type HealthCheckStore interface {
	CheckHealth(context.Context) (string, error)
}

type MongoHealthCheckStore struct {
	client *mongo.Client
}

func NewMongoHealthCheckStore(client *mongo.Client) *MongoHealthCheckStore {
	return &MongoHealthCheckStore{
		client: client,
	}
}

func (s *MongoHealthCheckStore) CheckHealth(ctx context.Context) (string, error) {
	cmd := bson.D{{Key: "serverStatus", Value: 1}}
	result := bson.D{}

	if err := s.client.Database(adminDB).RunCommand(ctx, cmd).Decode(&result); err != nil {
		return "", err
	}

	var version string
	for _, elem := range result {
		switch elem.Key {
		case "version":
			if val, ok := elem.Value.(string); ok {
				version = val
			}
		}
	}

	return version, nil
}
