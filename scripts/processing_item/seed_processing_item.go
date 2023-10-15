package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/johnson7543/ims/db"
	"github.com/johnson7543/ims/db/fixtures"
	"github.com/johnson7543/ims/types"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const processingItemColl = "processing_items"

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

	if err := client.Database(mongoDBName).Collection(processingItemColl).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	store := &db.Store{
		ProcessingItem: db.NewMongoProcessingItemStore(client),
	}

	workerId, err := primitive.ObjectIDFromHex("652374ca3feb85b51bf80ba1")
	if err != nil {
		log.Fatal(err)
	}
	sku, err := primitive.ObjectIDFromHex("6523bed323074e8c27b35ec1")
	if err != nil {
		log.Fatal(err)
	}
	for i := 1; i <= 5; i++ {
		processingItem := &types.ProcessingItem{
			Name:      fmt.Sprintf("Item_%d", i),
			Quantity:  i * 10,
			Price:     float64(i) * 5.0,
			WorkerID:  workerId,
			StartDate: time.Now().AddDate(0, 0, i),
			EndDate:   time.Now().AddDate(0, 0, i+5),
			SKU:       sku,
			Remarks:   fmt.Sprintf("Remarks for Item %d", i),
		}

		processingItem = fixtures.AddProcessingItem(store, processingItem)
		fmt.Println("Processing Item added ->", processingItem.Name)
	}
}
