package data

import (
	"os"

	"opa-lambda-database/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func Update(client *dynamodb.DynamoDB, user models.User) error {
	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("TABLE")),
		Item:      item,
	}

	_, err = client.PutItem(input)
	return err
}
