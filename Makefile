# Soulprint Backend Makefile
# ================================

# Variables
APP_NAME = soulprint-backend
MAIN_PATH = cmd/main.go
BUILD_DIR = build
BINARY_NAME = $(BUILD_DIR)/$(APP_NAME)

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
GORUN = $(GOCMD) run

# Colors for output
RED = \033[0;31m
GREEN = \033[0;32m
YELLOW = \033[0;33m
BLUE = \033[0;34m
PURPLE = \033[0;35m
CYAN = \033[0;36m
NC = \033[0m # No Color

.PHONY: help run build test clean deps install-deps start-mongo stop-mongo setup lint fmt vet check dev docker-build docker-run

# Default target
.DEFAULT_GOAL := help

## help: Display this help message
help:
	@echo "$(CYAN)Soulprint Backend - Available Commands:$(NC)"
	@echo ""
	@echo "$(GREEN)Development:$(NC)"
	@echo "  make run          - Run the application in development mode"
	@echo "  make dev          - Run with auto-reload (requires air)"
	@echo "  make build        - Build the application binary"
	@echo "  make test         - Run all tests"
	@echo "  make clean        - Clean build artifacts"
	@echo ""
	@echo "$(GREEN)Dependencies:$(NC)"
	@echo "  make deps         - Download and verify dependencies"
	@echo "  make install-deps - Install development dependencies"
	@echo "  make fmt          - Format Go code"
	@echo "  make lint         - Run linter (requires golint)"
	@echo "  make vet          - Run go vet"
	@echo "  make check        - Run fmt, vet, and lint"
	@echo ""
	@echo "$(GREEN)Database:$(NC)"
	@echo "  make start-mongo  - Start MongoDB service"
	@echo "  make stop-mongo   - Stop MongoDB service"
	@echo ""
	@echo "$(GREEN)Setup:$(NC)"
	@echo "  make setup        - Complete project setup"
	@echo "  make env          - Create .env file from template"
	@echo ""
	@echo "$(GREEN)Docker:$(NC)"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run application in Docker"

## run: Start the application in development mode
run: deps start-mongo
	@echo "$(GREEN)🚀 Starting Soulprint Backend...$(NC)"
	@echo "$(YELLOW)📍 Server will be available at http://localhost:8080$(NC)"
	@echo "$(YELLOW)📖 Health check: http://localhost:8080/health$(NC)"
	@echo "$(YELLOW)🔗 API docs: Check README.md for endpoints$(NC)"
	@echo ""
	$(GORUN) $(MAIN_PATH)

## dev: Run with auto-reload using air
dev: deps start-mongo
	@echo "$(GREEN)🚀 Starting Soulprint Backend with auto-reload...$(NC)"
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "$(RED)❌ Air not installed. Installing...$(NC)"; \
		go install github.com/cosmtrek/air@latest; \
		air; \
	fi

## build: Build the application binary
build: deps
	@echo "$(GREEN)🔨 Building $(APP_NAME)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✅ Binary created: $(BINARY_NAME)$(NC)"

## test: Run all tests
test: deps
	@echo "$(GREEN)🧪 Running tests...$(NC)"
	$(GOTEST) -v ./...

## clean: Clean build artifacts
clean:
	@echo "$(GREEN)🧹 Cleaning up...$(NC)"
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "$(GREEN)✅ Clean complete$(NC)"

## deps: Download and verify dependencies
deps:
	@echo "$(GREEN)📦 Downloading dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) verify
	$(GOMOD) tidy

## install-deps: Install development dependencies
install-deps:
	@echo "$(GREEN)🛠️  Installing development dependencies...$(NC)"
	@echo "Installing air for auto-reload..."
	@go install github.com/cosmtrek/air@latest
	@echo "Installing golangci-lint..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)✅ Development dependencies installed$(NC)"

## fmt: Format Go code
fmt:
	@echo "$(GREEN)🎨 Formatting code...$(NC)"
	$(GOCMD) fmt ./...

## lint: Run linter
lint:
	@echo "$(GREEN)🔍 Running linter...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)⚠️  golangci-lint not installed. Run 'make install-deps' first$(NC)"; \
	fi

## vet: Run go vet
vet:
	@echo "$(GREEN)🔍 Running go vet...$(NC)"
	$(GOCMD) vet ./...

## check: Run fmt, vet, and lint
check: fmt vet lint
	@echo "$(GREEN)✅ Code quality checks complete$(NC)"

## start-mongo: Start MongoDB service
start-mongo:
	@echo "$(GREEN)🍃 Starting MongoDB...$(NC)"
	@if brew services list | grep mongodb-community | grep started > /dev/null; then \
		echo "$(YELLOW)📄 MongoDB is already running$(NC)"; \
	else \
		brew services start mongodb/brew/mongodb-community; \
		echo "$(GREEN)✅ MongoDB started$(NC)"; \
	fi

## stop-mongo: Stop MongoDB service
stop-mongo:
	@echo "$(GREEN)🛑 Stopping MongoDB...$(NC)"
	@brew services stop mongodb/brew/mongodb-community
	@echo "$(GREEN)✅ MongoDB stopped$(NC)"

## setup: Complete project setup
setup: install-deps deps env start-mongo
	@echo "$(GREEN)🎉 Project setup complete!$(NC)"
	@echo ""
	@echo "$(CYAN)Next steps:$(NC)"
	@echo "1. Add your OpenAI API key to .env file"
	@echo "2. Run 'make run' to start the application"
	@echo "3. Import Postman collection for testing"

## env: Create .env file from template
env:
	@if [ ! -f .env ]; then \
		echo "$(GREEN)📄 Creating .env file...$(NC)"; \
		echo "PORT=8080" > .env; \
		echo "MONGODB_URI=mongodb://localhost:27017" >> .env; \
		echo "MONGODB_DATABASE=soulprint" >> .env; \
		echo "OPENAI_API_KEY=your_openai_api_key_here" >> .env; \
		echo "OPENAI_MODEL=gpt-3.5-turbo" >> .env; \
		echo "" >> .env; \
		echo "# Local Model Configuration" >> .env; \
		echo "USE_LOCAL_MODEL=true" >> .env; \
		echo "LOCAL_MODEL_URL=http://localhost:11434" >> .env; \
		echo "LOCAL_MODEL_NAME=llama3" >> .env; \
		echo "$(GREEN)✅ .env file created$(NC)"; \
		echo "$(YELLOW)⚠️  Remember to add your OpenAI API key!$(NC)"; \
	else \
		echo "$(YELLOW)📄 .env file already exists$(NC)"; \
	fi

## docker-build: Build Docker image
docker-build:
	@echo "$(GREEN)🐳 Building Docker image...$(NC)"
	@docker build -t $(APP_NAME):latest .

## docker-run: Run application in Docker
docker-run: docker-build
	@echo "$(GREEN)🐳 Running $(APP_NAME) in Docker...$(NC)"
	@docker run -p 8080:8080 --env-file .env $(APP_NAME):latest

# Additional utility targets
.PHONY: status logs restart

## status: Check application and database status
status:
	@echo "$(CYAN)📊 System Status:$(NC)"
	@echo ""
	@echo "$(GREEN)Application:$(NC)"
	@if curl -s http://localhost:8080/health > /dev/null 2>&1; then \
		echo "  ✅ API Server: Running"; \
		curl -s http://localhost:8080/health | grep -o '"status":"[^"]*"' || true; \
	else \
		echo "  ❌ API Server: Not running"; \
	fi
	@echo ""
	@echo "$(GREEN)Database:$(NC)"
	@if brew services list | grep mongodb-community | grep started > /dev/null; then \
		echo "  ✅ MongoDB: Running"; \
	else \
		echo "  ❌ MongoDB: Not running"; \
	fi

## logs: Show application logs (if running with systemd or similar)
logs:
	@echo "$(GREEN)📋 Application logs:$(NC)"
	@echo "For development, check terminal output where 'make run' was executed"

## restart: Restart the application (stop and start MongoDB)
restart: stop-mongo start-mongo
	@echo "$(GREEN)🔄 Services restarted$(NC)"
	@echo "Run 'make run' to start the application" 