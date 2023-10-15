build:
	@go build -o bin/api 

run: build
	@./bin/api

seed_user:
	@go run scripts/user/seed_user.go

seed_material:
	@go run scripts/material/seed_material.go

seed_worker:
	@go run scripts/worker/seed_worker.go

seed_product:
	@go run scripts/product/seed_product.go

seed_processing_item:
	@go run scripts/processing_item/seed_processing_item.go

seed_order:
	@go run scripts/order/seed_order.go

docker:
	echo "building docker file"
	@docker build -t api .
	echo "running API inside Docker container"
	@docker run -p 3000:3000 api

test:
	@go test -v ./...
