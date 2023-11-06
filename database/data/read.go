package data

import (
	"fmt"
	"os"

	"opa-lambda-database/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func Read(client *dynamodb.DynamoDB, userId string) (models.User, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE")),
		Key: map[string]*dynamodb.AttributeValue{
			"userId": {
				S: aws.String(userId),
			},
		},
	}

	result, err := client.GetItem(input)
	if err != nil {
		return models.User{}, err
	}
	if result.Item == nil {
		return models.User{}, fmt.Errorf("Cold not find '%s'", userId)
	}

	user := models.User{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &user)
	if err != nil {
		return models.User{}, fmt.Errorf("Failed to unmarshal Record, %v", err)
	}

	return user, nil
}
