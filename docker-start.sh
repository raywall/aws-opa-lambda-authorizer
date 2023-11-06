#!/bin/sh

NETWORK_NAME="local-api-network"
DYNAMO_IMAGE="amazon/dynamodb-local:latest"
DEBUG_CONTAINERS=""
WARM_CONTAINERS=""
REBUILD_IMAGES=""
REMOVE_IMAGES=0
PRIVATE_ENV=""

# Verifique os argumentos passados ao script
while [ "$#" -gt 0 ]; do
  case "$1" in
    --debug)
      DEBUG_CONTAINERS="--debug"
      ;;
    --warm)
      WARM_CONTAINERS="--warm-containers eager"
      ;;
    --rebuild)
      REBUILD_IMAGES="--no-cache"
      REMOVE_IMAGES=1
      ;;
    --private)
      PRIVATE_ENV="-t template-custom.yaml --skip-pull-image"
      ;;
    *)
      echo "Argumento desconhecido: $1"
      exit 1
      ;;
  esac
  shift
done

if [ "$PRIVATE_ENV" = "" ]; then
  PRIVATE_ENV="-t template.yaml"
  REMOVE_IMAGES=2
fi

docker ps -aq | xargs docker rm -f
docker network create $NETWORK_NAME

if [ "$REMOVE_IMAGES" = 1 ]; then
    # api lambda authorizer image
    docker rmi -f authorizer.ecr.aws/lambda/go:1.x
    docker build -f .docker/docker-authorizer -t authorizer.ecr.aws/lambda/go:1.x . $REBUILD_IMAGES

    # dynamodb local service
    docker rmi -f database.ecr.aws/dynamo/db:1.x
    docker build -f .docker/docker-dynamodb -t database.ecr.aws/dynamo/db:1.x . $REBUILD_IMAGES

    # database lambda function image
    docker rmi -f database.ecr.aws/lambda/go:1.x
    docker build -f .docker/docker-database -t database.ecr.aws/lambda/go:1.x . $REBUILD_IMAGES

    # hello lambda function image
    docker rmi -f hello.ecr.aws/lambda/go:1.x
    docker build -f .docker/docker-hello -t hello.ecr.aws/lambda/go:1.x . $REBUILD_IMAGES

    # token lambda function image
    docker rmi -f token.ecr.aws/lambda/go:1.x
    docker build -f .docker/docker-token -t token.ecr.aws/lambda/go:1.x . $REBUILD_IMAGES

    DYNAMO_IMAGE="database.ecr.aws/dynamo/db:1.x"
fi

if [ "$REMOVE_IMAGES" = 2 ]; then
  cd authorizer
  GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o main

  cd ../database
  GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o main

  cd ../hello
  GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o main

  cd ../token
  GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o main

  cd ..
  sam build $PRIVATE_ENV --docker-network $NETWORK_NAME --no-cached
fi

docker run -d -p 8000:8000 --network $NETWORK_NAME --hostname dynamodb --name dynamodb $DYNAMO_IMAGE
sleep 5

aws dynamodb create-table --endpoint-url http://localhost:8000 --table-name UserTable --attribute-definitions AttributeName="userId",AttributeType="S" --key-schema AttributeName="userId",KeyType=HASH --provisioned-throughput ReadCapacityUnits=2,WriteCapacityUnits=2

sam local start-api --docker-network $NETWORK_NAME $PRIVATE_ENV $DEBUG_CONTAINERS $WARM_CONTAINERS