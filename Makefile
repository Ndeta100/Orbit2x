# Makefile for orbit2x - Go application
.PHONY: help build run dev clean test templ deps install

# Variables
APP_NAME := orbit2x
BINARY_PATH := ./bin/$(APP_NAME)
GO_FILES := $(shell find . -name '*.go')
TEMPL_FILES := $(shell find . -name '*.templ')

# Colors
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m

## help: Show available commands
help:
	@echo "Available commands:"
	@echo "  make build    - Build the application"
	@echo "  make run      - Build and run the application"
	@echo "  make dev      - Run with hot reload (requires air)"
	@echo "  make templ    - Generate templ files"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make install  - Install dependencies"
	@echo "  make deps     - Update dependencies"

## install: Install dependencies and tools
install:
	@echo "$(YELLOW)Installing Go dependencies...$(NC)"
	go mod tidy
	@echo "$(YELLOW)Installing templ...$(NC)"
	go install github.com/a-h/templ/cmd/templ@latest
	@echo "$(YELLOW)Installing Node dependencies (for Tailwind)...$(NC)"
	npm install
	@echo "$(GREEN)Dependencies installed!$(NC)"

## deps: Update dependencies
deps:
	@echo "$(YELLOW)Updating Go dependencies...$(NC)"
	go get -u ./...
	go mod tidy
	npm update
	@echo "$(GREEN)Dependencies updated!$(NC)"

## templ: Generate templ files
templ:
	@echo "$(YELLOW)Generating templ files...$(NC)"
	templ generate
	@echo "$(GREEN)Templ files generated!$(NC)"

## build: Build the application
build: templ
	@echo "$(YELLOW)Building CSS assets...$(NC)"
	npm run build
	@echo "$(YELLOW)Building Go application...$(NC)"
	mkdir -p bin
	go build -o $(BINARY_PATH) .
	@echo "$(GREEN)Build complete!$(NC)"

## run: Build and run the application
run: build
	@echo "$(YELLOW)Starting $(APP_NAME)...$(NC)"
	$(BINARY_PATH)

## dev: Run with hot reload
dev: templ
	@echo "$(YELLOW)Building CSS assets...$(NC)"
	@echo "$(YELLOW)Starting development server...$(NC)"
	go run .

## test: Run tests
test:
	@echo "$(YELLOW)Running tests...$(NC)"
	go test -v ./...

## clean: Clean build artifacts
clean:
	@echo "$(YELLOW)Cleaning...$(NC)"
	rm -rf bin/
	rm -rf dist/
	go clean
	@echo "$(GREEN)Clean complete!$(NC)"

## format: Format Go code
format:
	@echo "$(YELLOW)Formatting code...$(NC)"
	go fmt ./...
	templ fmt .

# Default target
all: build