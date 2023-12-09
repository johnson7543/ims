build:
	@go build -o bin/ims-app 

run: build
	@./bin/ims-app

seed_user:
	@go run scripts/user/seed_user.go

seed_material:
	@go run scripts/material/seed_material.go

seed_worker:
	@go run scripts/worker/seed_worker.go

seed_customer:
	@go run scripts/customer/seed_customer.go
	
seed_seller:
	@go run scripts/seller/seed_seller.go

seed_buyer:
	@go run scripts/buyer/seed_buyer.go

seed_product:
	@go run scripts/product/seed_product.go

seed_processing_item:
	@go run scripts/processing_item/seed_processing_item.go

seed_order:
	@go run scripts/order/seed_order.go

docker:
	echo "building docker file"
	@docker build -t ims-app .
	echo "running ims-app inside Docker container"
	@docker run -p 3000:3000 ims-app

test:
	@go test -v ./...
