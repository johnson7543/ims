package fixtures

import (
	"context"
	"fmt"
	"log"

	"github.com/johnson7543/ims/db"
	"github.com/johnson7543/ims/types"
)

func AddOrder(store *db.Store, order *types.Order) *types.Order {
	insertedOrder, err := store.Order.InsertOrder(context.Background(), order)
	if err != nil {
		log.Fatal(err)
	}
	return insertedOrder
}

func AddProduct(store *db.Store, product *types.Product) *types.Product {
	insertedWorker, err := store.Product.InsertProduct(context.Background(), product)
	if err != nil {
		log.Fatal(err)
	}
	return insertedWorker
}

func AddProcessingItem(store *db.Store, processing_items *types.ProcessingItem) *types.ProcessingItem {
	insertedProcessingItem, err := store.ProcessingItem.InsertProcessingItem(context.Background(), processing_items)
	if err != nil {
		log.Fatal(err)
	}
	return insertedProcessingItem
}

func AddWorker(store *db.Store, worker *types.Worker) *types.Worker {
	insertedWorker, err := store.Worker.InsertWorker(context.Background(), worker)
	if err != nil {
		log.Fatal(err)
	}
	return insertedWorker
}

func AddCustomer(store *db.Store, customer *types.Customer) *types.Customer {
	insertedCustomer, err := store.Customer.InsertCustomer(context.Background(), customer)
	if err != nil {
		log.Fatal(err)
	}
	return insertedCustomer
}

func AddSeller(store *db.Store, seller *types.Seller) *types.Seller {
	insertedSeller, err := store.Seller.InsertSeller(context.Background(), seller)
	if err != nil {
		log.Fatal(err)
	}
	return insertedSeller
}

func AddMaterial(store *db.Store, material *types.Material) *types.Material {
	insertedMaterial, err := store.Material.InsertMaterial(context.Background(), material)
	if err != nil {
		log.Fatal(err)
	}
	return insertedMaterial
}

func AddUser(store *db.Store, fn, ln string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		FirstName: fn,
		LastName:  ln,
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin
	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}
