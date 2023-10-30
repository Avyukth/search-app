SHELL := /bin/bash
VERSION := 1.0
DOCKERFILE := Dockerfile
IMAGE_NAME := search-app
KIND            := kindest/node:v1.28.0
KIND_CLUSTER := "kind-deployment-cluster"
# =================================================================
# Run Commands
# =================================================================
run:
	go run cmd/server/main.go

test:
	go test ./... -count=1
	staticcheck -checks=all ./...

build-db:
	@echo "Building DB Docker Image..."
	docker-compose -f db.docker-compose.yml up -d

build: build-db
	@echo "Building Search App Docker Image..."
	DOCKER_BUILDKIT=1 docker build \
	--no-cache \
	-f $(DOCKERFILE) \
	-t $(IMAGE_NAME):$(VERSION) \
	--build-arg BUILD_REF=$(VERSION) \
	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
	.
build-search: 
	DOCKER_BUILDKIT=1 docker build \
	--no-cache \
	-f $(DOCKERFILE) \
	-t $(IMAGE_NAME):$(VERSION) \
	--build-arg BUILD_REF=$(VERSION) \
	--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
	.

docker-run:
	docker-compose up -d

docker-compose-up: build
	docker-compose up -d

docker-compose-down:
	docker-compose down

clean:
	@echo "Cleaning up..."
	@go clean
	@rm -f server

docker-log:
	docker logs --tail 100 -f $(IMAGE_NAME)

all: build docker-run


kind-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config deployment/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=search-system


kind-down:
	kind delete cluster --name $(KIND_CLUSTER)


kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces
