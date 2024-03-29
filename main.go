package main

import (
	"context"
	"log"
	"os"

	"github.com/johnson7543/ims/api"
	"github.com/johnson7543/ims/db"

	_ "github.com/johnson7543/ims/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoEndpoint).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	var (
		healthCheckStore    = db.NewMongoHealthCheckStore(client)
		userStore           = db.NewMongoUserStore(client)
		materialStore       = db.NewMongoMaterialStore(client)
		materialOrderStore  = db.NewMongoMaterialOrderStore(client)
		workerStore         = db.NewMongoWorkerStore(client)
		customerStore       = db.NewMongoCustomerStore(client)
		buyerStore          = db.NewMongoBuyerStore(client)
		sellerStore         = db.NewMongoSellerStore(client)
		porcessingItemStore = db.NewMongoProcessingItemStore(client)
		productStore        = db.NewMongoProductStore(client)
		orderStore          = db.NewMongoOrderStore(client)
		store               = &db.Store{
			HealthCheck:    healthCheckStore,
			User:           userStore,
			Material:       materialStore,
			MaterialOrder:  materialOrderStore,
			Worker:         workerStore,
			Customer:       customerStore,
			Buyer:          buyerStore,
			Seller:         sellerStore,
			ProcessingItem: porcessingItemStore,
			Product:        productStore,
			Order:          orderStore,
		}
		HealthCheckHandler    = api.NewHealthCheckHandler(store)
		authHandler           = api.NewAuthHandler(store)
		materialHandler       = api.NewMaterialHandler(store)
		materialOrderHandler  = api.NewMaterialOrderHandler(store)
		workerHandler         = api.NewWorkerHandler(store)
		customerHandler       = api.NewCustomerHandler(store)
		buyerHandler          = api.NewBuyerHandler(store)
		sellerHandler         = api.NewSellerHandler(store)
		processingItemHandler = api.NewProcessingItemHandler(store)
		productHandler        = api.NewProductHandler(store)
		orderHandler          = api.NewOrderHandler(store)
		app                   = fiber.New(config)
		homePage              = app.Group("/")
		healthCheck           = app.Group("/health")
		auth                  = app.Group("/api")
		apiv1                 = app.Group("/api/v1", api.JWTAuthentication(userStore))
	)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "X-Api-Token,Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowCredentials: true,
	}))

	app.Options("/*", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	})

	homePage.Get("/", HealthCheckHandler.HandleHealthCheck)
	healthCheck.Get("/", HealthCheckHandler.HandleHealthCheck)

	auth.Post("/auth", authHandler.HandleAuthenticate)

	apiv1.Get("/material", materialHandler.HandleGetMaterials)
	apiv1.Post("/material", materialHandler.HandleInsertMaterial)
	apiv1.Patch("/material/:id", materialHandler.HandleUpdateMaterial)
	apiv1.Delete("/material/:id", materialHandler.HandleDeleteMaterial)
	apiv1.Get("/material/colors", materialHandler.HandleGetMaterialColors)
	apiv1.Get("/material/types", materialHandler.HandleGetMaterialTypes)
	apiv1.Get("/material/sizes", materialHandler.HandleGetMaterialSizes)

	apiv1.Get("/materialOrder", materialOrderHandler.HandleGetMaterialOrders)
	apiv1.Post("/materialOrder", materialOrderHandler.HandleInsertMaterialOrder)
	apiv1.Patch("/materialOrder/:id", materialOrderHandler.HandleUpdateMaterialOrder)
	apiv1.Delete("/materialOrder/:id", materialOrderHandler.HandleDeleteMaterialOrder)
	apiv1.Post("/materialOrder/materialOrderItems/:id", materialOrderHandler.HandleInsertMaterialOrderItemsToOrder)

	apiv1.Get("/worker", workerHandler.HandleGetWorkers)
	apiv1.Post("/worker", workerHandler.HandleInsertWorker)
	apiv1.Patch("/worker/:id", workerHandler.HandleUpdateWorker)
	apiv1.Delete("/worker/:id", workerHandler.HandleDeleteWorker)

	apiv1.Get("/customer", customerHandler.HandleGetCustomers)
	apiv1.Post("/customer", customerHandler.HandleInsertCustomer)
	apiv1.Patch("/customer/:id", customerHandler.HandleUpdateCustomer)
	apiv1.Delete("/customer/:id", customerHandler.HandleDeleteCustomer)

	apiv1.Get("/buyer", buyerHandler.HandleGetBuyers)
	apiv1.Post("/buyer", buyerHandler.HandleInsertBuyer)
	apiv1.Patch("/buyer/:id", buyerHandler.HandleUpdateBuyer)
	apiv1.Delete("/buyer/:id", buyerHandler.HandleDeleteBuyer)

	apiv1.Get("/seller", sellerHandler.HandleGetSellers)
	apiv1.Post("/seller", sellerHandler.HandleInsertSeller)
	apiv1.Patch("/seller/:id", sellerHandler.HandleUpdateSeller)
	apiv1.Delete("/seller/:id", sellerHandler.HandleDeleteSeller)

	apiv1.Get("/processingItem", processingItemHandler.HandleGetProcessingItems)
	apiv1.Post("/processingItem", processingItemHandler.HandleInsertProcessingItem)
	apiv1.Patch("/processingItem/:id", processingItemHandler.HandleUpdateProcessingItem)
	apiv1.Delete("/processingItem/:id", processingItemHandler.HandleDeleteProcessingItem)

	apiv1.Get("/product", productHandler.HandleGetProducts)
	apiv1.Post("/product", productHandler.HandleInsertProduct)
	apiv1.Patch("/product/:id", productHandler.HandleUpdateProduct)
	apiv1.Delete("/product/:id", productHandler.HandleDeleteProduct)
	apiv1.Get("/product/colors", productHandler.HandleGetProductColors)
	apiv1.Get("/product/types", productHandler.HandleGetProductTypes)
	apiv1.Get("/product/sizes", productHandler.HandleGetProductSizes)

	apiv1.Get("/order", orderHandler.HandleGetOrders)
	apiv1.Post("/order", orderHandler.HandleInsertOrder)
	apiv1.Patch("/order/:id", orderHandler.HandleUpdateOrder)
	apiv1.Delete("/order/:id", orderHandler.HandleDeleteOrder)
	apiv1.Post("/order/orderItems/:id", orderHandler.HandleInsertOrderItemsToOrder)

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	app.Listen(listenAddr)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
