package main

import (
	"encoding/json"
	"fmt"
	"time"

	"opa-lambda-token/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("2cb6102a-e7c1-40a3-8c13-bc3270e95385")

func generate(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": "lista de features aqui",
		"exp":  time.Now().Add(time.Hour * 1).Unix(), // Expires in an hour
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("Token signing failed:", err)
		return "", err
	}

	return tokenString, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var user models.User
	err := json.Unmarshal([]byte(request.Body), &user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Failed to generate token: %v", err.Error()),
			StatusCode: 501,
		}, nil
	}

	token, err := generate(user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Failed to generate token: %v", err.Error()),
			StatusCode: 501,
		}, nil
	}

	return events.APIGatewayProxyResponse{Body: token, StatusCode: 200}, nil
}

func main() {
	lambda.Start(handler)
}
