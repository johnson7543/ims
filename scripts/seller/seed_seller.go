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

const sellerColl = "sellers"

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

	if err := client.Database(mongoDBName).Collection(sellerColl).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	store := &db.Store{
		Seller: db.NewMongoSellerStore(client),
	}

	for i := 1; i <= 5; i++ {
		seller := &types.Seller{
			Company:     fmt.Sprintf("SellerCompany_%d", i),
			Name:        fmt.Sprintf("SellerName_%d", i),
			Phone:       fmt.Sprintf("SellerPhone_%d", i),
			Address:     fmt.Sprintf("SellerAddress_%d", i),
			TaxIdNumber: fmt.Sprintf("SellerTaxId_%d", i),
		}

		seller = fixtures.AddSeller(store, seller)
		fmt.Println("Seller added ->", seller.Name)
	}
}
