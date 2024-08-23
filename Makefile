.PHONY: vendor serve fixtures docker-login docker-build docker-tag docker-push

# Clean and update dependencies.
vendor:
	go mod tidy && go mod vendor

# Start the application server.
serve:
	go run cmd/server/main.go

# Reset the database to a known state.
fixtures:
	go run cmd/fixtures/main.go

docker-login:
	aws ecr get-login-password | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.us-east-1.amazonaws.com

docker-build:
	docker build -t kitchens-api .

docker-tag:
	docker tag kitchens-api:latest ${AWS_ACCOUNT_ID}.dkr.ecr.us-east-1.amazonaws.com/kitchens_api_production:latest

docker-push:
	docker push ${AWS_ACCOUNT_ID}.dkr.ecr.us-east-1.amazonaws.com/kitchens_api_production:latest