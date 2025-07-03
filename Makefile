# TLD Scanner Makefile

# Variables
BINARY_NAME=tldscanner
VERSION=2.0.0
BUILD_DIR=build
MAIN_FILE=tldscanner.go

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION) -s -w"
BUILD_FLAGS=-trimpath

.PHONY: all build clean test deps run install help

# Default target
all: clean deps build

# Build the application
build:
	@echo "Building $(BINARY_NAME) v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
build-all: clean deps
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	
	# Linux AMD64
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	
	# Linux ARM64
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_FILE)
	
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)
	
	# macOS ARM64 (M1/M2)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_FILE)
	
	# Windows AMD64
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_FILE)
	
	# FreeBSD AMD64
	GOOS=freebsd GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-freebsd-amd64 $(MAIN_FILE)
	
	@echo "Multi-platform build completed!"
	@ls -la $(BUILD_DIR)/

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) tidy
	$(GOMOD) download

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run the application with example parameters
run: build
	@echo "Running $(BINARY_NAME) with example.com..."
	./$(BUILD_DIR)/$(BINARY_NAME) -d example.com -w wordlist.txt -v

# Run with custom domain
run-domain: build
	@if [ -z "$(DOMAIN)" ]; then \
		echo "Usage: make run-domain DOMAIN=example.com"; \
		exit 1; \
	fi
	./$(BUILD_DIR)/$(BINARY_NAME) -d $(DOMAIN) -w wordlist.txt -v

# Install to system PATH
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "Installation completed!"

# Uninstall from system
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstallation completed!"

# Create release packages
release: build-all
	@echo "Creating release packages..."
	@mkdir -p $(BUILD_DIR)/releases
	
	# Linux AMD64
	tar -czf $(BUILD_DIR)/releases/$(BINARY_NAME)-v$(VERSION)-linux-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-linux-amd64 -C .. wordlist.txt README.md
	
	# Linux ARM64
	tar -czf $(BUILD_DIR)/releases/$(BINARY_NAME)-v$(VERSION)-linux-arm64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-linux-arm64 -C .. wordlist.txt README.md
	
	# macOS AMD64
	tar -czf $(BUILD_DIR)/releases/$(BINARY_NAME)-v$(VERSION)-darwin-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-darwin-amd64 -C .. wordlist.txt README.md
	
	# macOS ARM64
	tar -czf $(BUILD_DIR)/releases/$(BINARY_NAME)-v$(VERSION)-darwin-arm64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-darwin-arm64 -C .. wordlist.txt README.md
	
	# Windows AMD64
	zip -j $(BUILD_DIR)/releases/$(BINARY_NAME)-v$(VERSION)-windows-amd64.zip $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe wordlist.txt README.md
	
	# FreeBSD AMD64
	tar -czf $(BUILD_DIR)/releases/$(BINARY_NAME)-v$(VERSION)-freebsd-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-freebsd-amd64 -C .. wordlist.txt README.md
	
	@echo "Release packages created in $(BUILD_DIR)/releases/"
	@ls -la $(BUILD_DIR)/releases/

# Generate checksums for releases
checksums: release
	@echo "Generating checksums..."
	cd $(BUILD_DIR)/releases && sha256sum * > checksums.txt
	@echo "Checksums generated!"

# Format code
fmt:
	@echo "Formatting code..."
	$(GOCMD) fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint not installed. Installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run

# Security check
security:
	@echo "Running security checks..."
	@if ! command -v gosec &> /dev/null; then \
		echo "gosec not installed. Installing..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
	fi
	gosec ./...

# Benchmark
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Profile
profile: build
	@echo "Running with profiling..."
	./$(BUILD_DIR)/$(BINARY_NAME) -d example.com -w wordlist.txt -cpuprofile=cpu.prof -memprofile=mem.prof

# Development mode with file watching
dev:
	@echo "Starting development mode..."
	@if ! command -v air &> /dev/null; then \
		echo "air not installed. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	air

# Create sample wordlist
sample-wordlist:
	@echo "Creating sample wordlist..."
	@echo -e "com\nnet\norg\nedu\ngov\nco\nio\nme\ntv\ncc\nbiz\ninfo" > sample_wordlist.txt
	@echo "Sample wordlist created: sample_wordlist.txt"

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) .

# Docker run
docker-run:
	@if [ -z "$(DOMAIN)" ]; then \
		echo "Usage: make docker-run DOMAIN=example.com"; \
		exit 1; \
	fi
	docker run --rm -it $(BINARY_NAME):$(VERSION) -d $(DOMAIN)

# Help
help:
	@echo "TLD Scanner Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build          Build the application"
	@echo "  build-all      Build for multiple platforms"
	@echo "  deps           Install dependencies"
	@echo "  clean          Clean build artifacts"
	@echo "  test           Run tests"
	@echo "  run            Run with example domain"
	@echo "  run-domain     Run with custom domain (make run-domain DOMAIN=example.com)"
	@echo "  install        Install to system PATH"
	@echo "  uninstall      Uninstall from system"
	@echo "  release        Create release packages"
	@echo "  checksums      Generate checksums for releases"
	@echo "  fmt            Format code"
	@echo "  lint           Lint code"
	@echo "  security       Run security checks"
	@echo "  bench          Run benchmarks"
	@echo "  profile        Run with profiling"
	@echo "  dev            Start development mode with file watching"
	@echo "  sample-wordlist Create sample wordlist"
	@echo "  docker-build   Build Docker image"
	@echo "  docker-run     Run Docker container (make docker-run DOMAIN=example.com)"
	@echo "  help           Show this help message"
