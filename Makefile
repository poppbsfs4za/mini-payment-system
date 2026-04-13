APP_NAME=arise-assignment

run:
	go run ./cmd/server

test:
	go test ./...

build:
	go build -o bin/$(APP_NAME) ./cmd/server

swagger:
	swag init -g cmd/server/main.go -o docs

openapi:
	@echo "Open docs/openapi.yaml in Swagger Editor or import it into Postman"
