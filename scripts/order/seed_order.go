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

const orderColl = "orders"

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

	if err := client.Database(mongoDBName).Collection(orderColl).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	store := &db.Store{
		Order: db.NewMongoOrderStore(client),
	}

	productIdStrings := []string{
		"653f2b14afc3a62643af2ab4",
		"653f2b14afc3a62643af2ab5",
		"653f2b14afc3a62643af2ab6",
		"653f2b14afc3a62643af2ab7",
		"653f2b15afc3a62643af2ab8",
	}
	productIds := make([]primitive.ObjectID, len(productIdStrings))
	for i, idStr := range productIdStrings {
		sku, err := primitive.ObjectIDFromHex(idStr)
		if err != nil {
			fmt.Printf("Error converting hex to ObjectID for %s: %v\n", idStr, err)
			return
		}
		productIds[i] = sku
	}

	for i := 1; i <= 5; i++ {

		orderItems := make([]types.OrderItem, 0)

		for j := 0; j < 3; j++ {
			orderProduct := &types.OrderProduct{
				SKU:       productIds[j],
				UnitPrice: 10 + (10 * float64(j+1)),
			}

			orderItem := &types.OrderItem{
				Product:    *orderProduct,
				Quantity:   i,
				TotalPrice: orderProduct.UnitPrice * float64(j+1),
			}

			orderItems = append(orderItems, *orderItem)
		}

		order := &types.Order{
			CustomerID:      primitive.NewObjectID(),
			OrderDate:       time.Now(),
			DeliveryDate:    time.Now().Add(time.Hour * 24 * 7),
			PaymentDate:     time.Now().Add(time.Hour * 24 * 14),
			TotalAmount:     500.0,
			Status:          "Delivered",
			ShippingAddress: "Taipei City",
			OrderItems:      orderItems,
		}

		order = fixtures.AddOrder(store, order)
		fmt.Println("Product added ->", order.ID)
	}
}
