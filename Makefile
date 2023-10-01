.PHONY: build test run clean docker-build docker-run docker-compose-up docker-compose-down

# Variables
BINARY_NAME=search-api
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

# Build Docker Image
docker-build:
	@echo "Building Docker Image..."
	@docker build -t search-api .

# Run Docker Container
docker-run:
	@echo "Running Docker Container..."
	@docker run -p 8080:8080 search-api

# Start services using Docker Compose
docker-compose-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d

# Stop services using Docker Compose
docker-compose-down:
	@echo "Stopping services with Docker Compose..."
	@docker-compose down

# View logs of services started with Docker Compose
docker-compose-logs:
	@echo "Viewing logs of services started with Docker Compose..."
	@docker-compose logs -f
