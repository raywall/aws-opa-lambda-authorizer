package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"go-opa-lambda-authorizer/models"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dgrijalva/jwt-go"
	"github.com/open-policy-agent/opa/rego"
)

var (
	role      string
	opaPolicy *rego.Rego
	userId    string

	//go:embed data/*
	fs embed.FS

	secretKey = []byte("2cb6102a-e7c1-40a3-8c13-bc3270e95385")
)

func init() {
	regoContent, err := fs.ReadFile("data/rules.rego")
	if err != nil {
		log.Fatal("Failed to read Rego policy file:", err)
	}

	opaPolicy = rego.New(
		rego.Query("data.rules.allow"),
		rego.Module("data/rules.rego", string(regoContent)),
	)
}

func main() {
	lambda.Start(authorize)
}

func validate(bearerToken string) (*jwt.Token, error) {
	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("Invalid signature method: %s", token.Header["alg"])
		}

		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Unauthorized: Token verification failed: %s (token received: %s)", err, bearerToken)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		role = claims["role"].(string)
		//userID := claims["sub"].(string)
		expiration := claims["exp"].(float64)
		expirationTime := time.Unix(int64(expiration), 0)

		if time.Now().After(expirationTime) {
			return nil, errors.New("Unauthorized: Token expired")
		}
	}

	return token, nil
}

func authorize(ctx context.Context, request events.APIGatewayCustomAuthorizerRequest) (events.APIGatewayCustomAuthorizerResponse, error) {
	dataContent, err := fs.ReadFile("data/data.json")
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, fmt.Errorf("Unauthorized: Failed to read data file: %s", err)
	}

	var data models.Data
	err = json.Unmarshal(dataContent, &data)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, fmt.Errorf("Unauthorized: Failed to unmarshal data: %s", err)
	}

	token, err := validate(strings.TrimPrefix(request.AuthorizationToken, "Bearer "))
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, fmt.Errorf("Unauthorized: Token verification failed: %s", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		role = claims["role"].(string)
		//userID := claims["sub"].(string)
		expiration := claims["exp"].(float64)
		expirationTime := time.Unix(int64(expiration), 0)

		if time.Now().After(expirationTime) {
			return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized: Token expired")
		}
	}

	methodArn := models.MethodARN{}
	if fail := methodArn.Load(request.MethodArn); fail != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, fmt.Errorf("Unauthorized: Faield to unmarshal method arn: %s", fail)
	}

	input := map[string]interface{}{
		"method":   methodArn.Method,
		"path":     methodArn.Path,
		"userType": "admin",
		"role":     role,
	}

	pq, err := opaPolicy.PrepareForEval(ctx)
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, fmt.Errorf("Unauthorized: Failed to prepare policy for eval: %s", err)
	}

	rs, err := pq.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return events.APIGatewayCustomAuthorizerResponse{}, fmt.Errorf("Unauthorized: Failed to evaluate policy: %s", err)
	}

	if rs.Allowed() {
		return generateAllowPolicyResponse(methodArn.Method), nil
	}

	return events.APIGatewayCustomAuthorizerResponse{}, errors.New("Unauthorized")
}

func generateAllowPolicyResponse(method string) events.APIGatewayCustomAuthorizerResponse {
	return events.APIGatewayCustomAuthorizerResponse{
		PrincipalID:    "user",
		PolicyDocument: generatePolicyDocument(method),
	}
}

func generatePolicyDocument(method string) events.APIGatewayCustomAuthorizerPolicy {
	return events.APIGatewayCustomAuthorizerPolicy{
		Version: "2012-10-17",
		Statement: []events.IAMPolicyStatement{
			{
				Action:   []string{"execute-api:Invoke"},
				Effect:   "Allow",
				Resource: []string{"*"},
			},
		},
	}
}
