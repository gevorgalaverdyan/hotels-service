APP_NAME=app
BUILD_DIR=./bin

fmt:
	gofmt -w .

dev:
	go run main.go

build:
	mkdir -p $(BUILD_DIR)
	go build -v -o $(BUILD_DIR)/$(APP_NAME) .
