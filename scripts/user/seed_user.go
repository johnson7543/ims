package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/johnson7543/ims/api"
	"github.com/johnson7543/ims/db"
	"github.com/johnson7543/ims/db/fixtures"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const userColl = "users"

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

	if err := client.Database(mongoDBName).Collection(userColl).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	store := &db.Store{
		User: db.NewMongoUserStore(client),
	}

	// user := fixtures.AddUser(store, "johnsonwang", "test", false)
	// fmt.Println("johnsonwang ->", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println("admin ->", api.CreateTokenFromUser(admin))

}
