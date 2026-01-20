.PHONY: run build test clean install

# Run the application
run:
	go run cmd/server/main.go

# Build the application
build:
	go build -o bin/biller-server cmd/server/main.go

# Install dependencies
install:
	go mod download
	go mod tidy

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Format code
fmt:
	go fmt ./...

# Run linter
lint:
	golangci-lint run

# Create .env from example
env:
	cp .env.example .env
