.PHONY: all build test lint clean generate-mocks

# Default target
all: test build

# Build the application
build:
	go build -o bin/cc ./cmd/cc

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
test-cover:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Run unit tests only
test-unit:
	go test -v ./internal/... ./pkg/...

# Run integration tests only
test-integration:
	go test -v ./test/integration/...

# Run end-to-end tests only
test-e2e:
	go test -v ./test/e2e/...

# Run linters
lint:
	go vet ./...
	if command -v staticcheck > /dev/null; then \
		staticcheck ./...; \
	else \
		echo "staticcheck not installed, skipping..."; \
	fi

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Generate mocks
generate-mocks:
	go install go.uber.org/mock/mockgen@v0.4.0
	go run tools/generate_mocks.go

# Format code
fmt:
	go fmt ./...