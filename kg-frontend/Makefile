.PHONY: build test run clean docker-build docker-run

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=kg-frontend
SRC_DIR=src
PUBLIC_DIR=public

# Docker parameters
DOCKER_IMAGE=kg-frontend
DOCKER_TAG=latest

all: test build

build:
	cd $(SRC_DIR) && $(GOBUILD) -o ../$(BINARY_NAME) -v

test:
	cd $(SRC_DIR) && $(GOTEST) -v ./...

run: build
	./$(BINARY_NAME)

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

deps:
	cd $(SRC_DIR) && $(GOMOD) download

docker-build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

help:
	@echo "make - Build and test the application"
	@echo "make build - Build the application"
	@echo "make test - Run tests"
	@echo "make run - Run the application"
	@echo "make clean - Clean build files"
	@echo "make deps - Download dependencies"
	@echo "make docker-build - Build Docker image"
	@echo "make docker-run - Run Docker container" 