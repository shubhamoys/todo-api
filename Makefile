# Makefile for Todo API

# Variables
APP_NAME=todo-api
MAIN_PATH=cmd/main.go
BUILD_DIR=build

# Build the application
build:
	go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_PATH)

# Run the application
run:
	go run $(MAIN_PATH)

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)

# Install dependencies
deps:
	go mod download

# Seed the superadmin user
seed-admin:
	go run cmd/seed_superadmin/seed_superadmin.go

# Start development environment with hot reload
dev:
	air -c .air.toml

# Help command to list available targets
help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make deps         - Install dependencies"
	@echo "  make seed-admin   - Seed the superadmin user"
	@echo "  make dev          - Start development with hot reload"

.PHONY: build run clean deps seed-admin dev help
