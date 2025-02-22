fmt:
	gofmt -w .

dev:
	go run main.go

build:
	go build -o bin/ ./...