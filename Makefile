.PHONY: vendor serve fixtures

# Clean and update dependencies.
vendor:
	go mod tidy && go mod vendor

# Start the application server.
serve:
	go run cmd/server/main.go

# Reset the database to a known state.
fixtures:
	go run cmd/fixtures/main.go