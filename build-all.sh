#!/bin/sh

cd authorizer
echo "Building authorizer ..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o main

cd ../database
echo "Building database ..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o main

cd ../hello
echo "Building hello ..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o main

cd ../token
echo "Building token ..."
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o main

cd ..
echo "Done!"