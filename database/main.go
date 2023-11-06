package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"opa-lambda-database/data"
	"opa-lambda-database/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
	client *dynamodb.DynamoDB
)

func init() {
	sess := session.Must(session.NewSession())
	config := aws.NewConfig().WithRegion(os.Getenv("AWS_REGION"))

	if os.Getenv("AWSENV") == "AWS_SAM_LOCAL" {
		config = config.WithEndpoint(os.Getenv("ENDPOINT"))
	}

	client = dynamodb.New(sess, config)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == "POST" {
		var user models.User
		err := json.Unmarshal([]byte(request.Body), &user)

		if err != nil {
			log.Println("Failed to unmarshal user:", err)

			return events.APIGatewayProxyResponse{
				Body:       fmt.Sprintf("Failed to read record: %s", err.Error()),
				StatusCode: 400,
			}, nil
		}

		err = data.Update(client, user)
		if err != nil {
			log.Println("Failed to insert record:", err)

			return events.APIGatewayProxyResponse{
				Body:       fmt.Sprintf("Failed to read record: %s", err.Error()),
				StatusCode: 500,
			}, nil
		}

		return events.APIGatewayProxyResponse{StatusCode: 201}, nil

	} else if request.HTTPMethod == "GET" {
		userId := request.PathParameters["userId"]
		record, err := data.Read(client, userId)

		if err != nil {
			log.Println("Failed to read record:", err)
			return events.APIGatewayProxyResponse{StatusCode: 500}, nil
		}

		response, _ := json.Marshal(record)
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    map[string]string{"Content-Type": "application/json"},
			Body:       string(response),
		}, nil
	}

	return events.APIGatewayProxyResponse{StatusCode: 405}, nil
}

func main() {
	lambda.Start(handler)
}
