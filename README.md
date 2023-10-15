# IMS

## Project environment variables

.env file
```
HTTP_LISTEN_ADDRESS=:3333
JWT_SECRET=
MONGO_DB_NAME=
MONGO_DB_URL=

```

## Project outline

- Users -> for staffs to access IMS
- Authentication and authorization -> JWT tokens
- Material -> CRUD API -> JSON
- Worker -> CRUD API -> JSON
- Processing item -> CRUD API -> JSON
- Product -> CRUD API -> JSON
- Order -> CRUD API -> JSON
- Scripts -> database management -> seeding

## Resources

### Mongodb driver

Documentation

```
https://mongodb.com/docs/drivers/go/current/quick-start
```

Installing mongodb client

```
go get go.mongodb.org/mongo-driver/mongo
```

### gofiber

Documentation

```
https://gofiber.io
```

Installing gofiber

```
go get github.com/gofiber/fiber/v2
```

## Docker

### Installing mongodb as a Docker container

```
docker run --name mongodb -d mongo:latest
```

## Swaggo

### Installing swaggo

```
go get -u github.com/swaggo/swag/cmd/swag
```

### Generating swagger docs

```
swag init --parseDependency --parseInternal
```

### Start/Stop the server

```bash
sudo systemctl daemon-reload
sudo systemctl restart ims-app
sudo systemctl status ims-app

sudo systemctl stop ims-app
sudo pkill -f ims-app
```