.PHONY: mod clean build start d-start d-stop dynamodb-admin create-table insert-data help

# Show usage if no option
.DEFAULT_GOAL := help

mod: ## Get dependencies
	go get -u ./...

clean: ## Remove binary
	rm -rf ./authorizer/main
	rm -rf ./database/main
	rm -rf ./hello/main
	rm -rf ./token/main
	
build: ## Build binary
	GOOS=linux GOARCH=x86_64 go build -o authorizer/main ./authorizer
	GOOS=linux GOARCH=x86_64 go build -o database/main ./database
	GOOS=linux GOARCH=x86_64 go build -o hello/main ./hello
	GOOS=linux GOARCH=x86_64 go build -o token/main ./token

start: build ## Start Lambda and API Gateway on localhost
	sam local start-api --env-vars db/env.json # or start-lambda

d-start: ## Boot DynamoDB local
	docker-compose up -d

d-stop: ## Treminate DynamoDB local
	docker-compose down

pre-admin:
	@if [ -z `which dynamodb-admin 2> /dev/null` ]; then \
		echo "Need to install dynamodb-admin, execute \"npm install dynamodb-admin -g\"";\
		exit 1;\
	fi

dynamodb-admin: pre-admin ## Start DaynamoDB GUI
	DYNAMO_ENDPOINT=http://localhost:18000 dynamodb-admin

create-table:
	aws dynamodb create-table --cli-input-json file://db/create_user_table.json --endpoint-url http://localhost:18000

insert-data:
	aws dynamodb batch-write-item --request-items file://db/batch_data.json --endpoint-url http://localhost:18000

help: ## Show options
	@grep -E '^[a-zA-Z_\.-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
