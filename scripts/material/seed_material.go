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

const materialColl = "materials"

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

	if err := client.Database(mongoDBName).Collection(materialColl).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	store := &db.Store{
		Material: db.NewMongoMaterialStore(client),
	}

	for i := 1; i <= 5; i++ {
		material := &types.Material{
			Name:     fmt.Sprintf("Material_%d", i),
			Color:    fmt.Sprintf("Color_%d", i), // You can replace this with actual color data
			Size:     fmt.Sprintf("Size_%d", i),  // You can replace this with actual size data
			Quantity: fmt.Sprintf("%d", i*10),
			Remarks:  fmt.Sprintf("Remarks for Material %d", i),
		}

		material = fixtures.AddMaterial(store, material)
		fmt.Println("Material added ->", material.Name)

	}

}
