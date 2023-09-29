.PHONY: build test run clean

# Variables
BINARY_NAME=server
PKG_NAME=cmd/server

# Build the project
build:
	@echo "Building..."
	@go build -o $(BINARY_NAME) ./$(PKG_NAME)

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Run the project
run: build
	@echo "Running..."
	@./$(BINARY_NAME)

# Clean up
clean:
	@echo "Cleaning up..."
	@go clean
	@rm -f $(BINARY_NAME)
