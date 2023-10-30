package db

const MongoDBNameEnvName = "MONGO_DB_NAME"

type Pagination struct {
	Limit int64
	Page  int64
}

type Store struct {
	HealthCheck    HealthCheckStore
	User           UserStore
	Material       MaterialStore
	MaterialOrder  MaterialOrderStore
	Worker         WorkerStore
	Product        ProductStore
	ProcessingItem ProcessingItemStore
	Order          OrderStore
}
