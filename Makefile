test:
	golangci-lint run ./...
	go test ./... -v -cover

mock:
	mockgen -source=./internal/repository/manager.go -destination=./internal/repository/mock/mock.go
	mockgen -source=./internal/service/manager.go -destination=./internal/service/mock/mock.go
	mockgen -source=./pkg/clients/binance/binance.go -destination=./pkg/clients/binance/mock/mock.go

run:
	go run ./cmd/main.go -config_path=local.env

swag:
	swag init -g ./cmd/main.go -o docs