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

sync:
	go run cmd/sync/main.go

# Log into Docker
docker-login:
	aws ecr get-login-password | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.us-east-1.amazonaws.com

# Build the image
docker-build:
	docker build -t kitchens-api:${API_VERSION} .

# Tag the image with Latest
docker-tag:
	docker tag kitchens-api:latest ${AWS_ACCOUNT_ID}.dkr.ecr.us-east-1.amazonaws.com/kitchens_api_production:latest

# Version
docker-version:
	docker tag kitchens-api:${API_VERSION} ${AWS_ACCOUNT_ID}.dkr.ecr.us-east-1.amazonaws.com/kitchens_api_production:${API_VERSION}

# Push the Docker image to ECR
docker-push:
	docker push ${AWS_ACCOUNT_ID}.dkr.ecr.us-east-1.amazonaws.com/kitchens_api_production:latest

docker-version-push:
	docker push ${AWS_ACCOUNT_ID}.dkr.ecr.us-east-1.amazonaws.com/kitchens_api_production:${API_VERSION}

docker-release: docker-build docker-version docker-version-push