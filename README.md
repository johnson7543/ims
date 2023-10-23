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

### Swaggo

#### Installing swaggo

```
go get -u github.com/swaggo/swag/cmd/swag
```

#### Generating swagger docs

```
swag init --parseDependency --parseInternal
```

### Run the Server

#### Start/Stop the server

```bash
sudo systemctl daemon-reload
sudo systemctl restart ims-app
sudo systemctl status ims-app

sudo systemctl stop ims-app
sudo pkill -f ims-app
```

#### Run with Docker

```bash
# login in docker with AWS
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin 248679804578.dkr.ecr.ap-northeast-1.amazonaws.com

docker build -t ims-ecs:v1 .
docker tag ims-ecs:v1 248679804578.dkr.ecr.ap-northeast-1.amazonaws.com/ims:latest
docker push 248679804578.dkr.ecr.ap-northeast-1.amazonaws.com/ims:latest

# run in local
docker run --name ims-local -p 8080:8080 ims-ecs:v1
```

#### Circle CI reference

<https://circleci.com/blog/use-circleci-orbs-to-build-test-and-deploy-a-simple-go-application-to-aws-ecs/>
<https://circleci.com/developer/orbs/orb/circleci/aws-ecr#usage-simple_build_and_push>
<https://circleci.com/developer/orbs/orb/circleci/aws-ecs#usage-deploy_service_update>
<https://circleci.com/docs/openid-connect-tokens/>

#### Some AWS cmd

```bash
aws ecs register-task-definition --cli-input-json file://task-definition.json
aws ecs create-service --cli-input-json file://ecs-service.json
```
