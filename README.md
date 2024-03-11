# IMS

## Project environment variables setup

create .env file at root

```bash
HTTP_LISTEN_ADDRESS=:8080
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

```bash
go get -u github.com/swaggo/swag/cmd/swag
```

#### Generating Swagger docs

```bash
swag init --parseDependency --parseInternal
```

##### Swagger link

<http://localhost:8080/swagger/index.html#/>

### Run the App

#### Run with Docker

```bash
# login in docker with AWS
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin 471112724062.dkr.ecr.ap-northeast-1.amazonaws.com

# build and push to aws ecr
# docker build -t ims-ecs:v1 .
docker build -t ims-ecs:v1 --build-arg ENV_MONGO_DB_PASSWORD=root . 
docker tag ims-ecs:v1 471112724062.dkr.ecr.ap-northeast-1.amazonaws.com/ims:latest
docker push 471112724062.dkr.ecr.ap-northeast-1.amazonaws.com/ims:latest

# build and run in local
docker build -t ims-ecs:v1 --build-arg ENV_MONGO_DB_PASSWORD=root . 
docker run --name ims-local -p 8080:8080 ims-ecs:v1
```

#### Circle CI reference

<https://circleci.com/blog/use-circleci-orbs-to-build-test-and-deploy-a-simple-go-application-to-aws-ecs/>
<https://github.com/CircleCI-Public/aws-ecr-orb/blob/master/src/commands/build_and_push_image.yml>
<https://github.com/CircleCI-Public/aws-ecs-orb/blob/master/src/commands/update_service.yml>
<https://circleci.com/developer/orbs/orb/circleci/aws-ecr#usage-simple_build_and_push>
<https://circleci.com/developer/orbs/orb/circleci/aws-ecs#usage-deploy_service_update>
<https://circleci.com/docs/openid-connect-tokens/>

#### Some AWS cmd

```bash
aws ecs register-task-definition --cli-input-json file://task-definition.json
aws ecs create-service --cli-input-json file://ecs-service.json
```
