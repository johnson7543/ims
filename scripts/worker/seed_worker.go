package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/johnson7543/ims/db"
	"github.com/johnson7543/ims/db/fixtures"
	"github.com/johnson7543/ims/types"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const workerColl = "workers"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		ctx           = context.Background()
		mongoEndpoint = os.Getenv("MONGO_DB_URL")
		mongoDBName   = os.Getenv("MONGO_DB_NAME")
	)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoEndpoint).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	if err := client.Database(mongoDBName).Collection(workerColl).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	store := &db.Store{
		Worker: db.NewMongoWorkerStore(client),
	}

	for i := 1; i <= 5; i++ {
		worker := &types.Worker{
			Company:     fmt.Sprintf("Company_%d", i),
			Name:        fmt.Sprintf("Name_%d", i),
			Phone:       fmt.Sprintf("Phone_%d", i),
			Address:     fmt.Sprintf("Address_%d", i),
			TaxIdNumber: fmt.Sprintf("TaxId_%d", i),
		}

		worker = fixtures.AddWorker(store, worker)
		fmt.Println("Worker added ->", worker.Name)
	}

}
