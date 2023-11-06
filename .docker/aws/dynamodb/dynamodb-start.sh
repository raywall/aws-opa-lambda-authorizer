#!/bin/sh -e

java -Djava.library.path=/dynamodb/DynamoDBLocal_lib -jar /dynamodb/DynamoDBLocal.jar -sharedDb -port 8000 -dbPath /dynamodb/tables/ &
sleep 5

aws dynamodb create-table --endpoint-url http://localhost:8000 --table-name UserTable --attribute-definitions AttributeName="userId",AttributeType="S" --key-schema AttributeName="userId",KeyType=HASH --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1

exit 0