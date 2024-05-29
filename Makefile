.PHONY: vendor serve

vendor:
	go mod tidy && go mod vendor

serve:
	go run cmd/server/main.go