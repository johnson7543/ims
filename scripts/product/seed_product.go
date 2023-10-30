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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const productColl = "products"

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

	if err := client.Database(mongoDBName).Collection(productColl).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	store := &db.Store{
		Product: db.NewMongoProductStore(client),
	}

	for i := 1; i <= 5; i++ {
		product := &types.Product{
			SKU:      fmt.Sprintf("SKU_%d", i),
			Name:     fmt.Sprintf("Product_%d", i),
			Material: fmt.Sprintf("Material_%d", i),
			Color:    fmt.Sprintf("Color_%d", i),
			Size:     fmt.Sprintf("Size_%d", i),
			Quantity: i * 20,
			Price:    float64(i) * 10.0,
			Date:     time.Now().AddDate(0, 0, i),
			Remark:   fmt.Sprintf("Remark for Product %d", i),
		}

		product = fixtures.AddProduct(store, product)
		fmt.Println("Product added ->", product.Name)
	}
}
