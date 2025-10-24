# Kinde Go SDK Makefile

.PHONY: help generate build test clean install-tools

# Default target
help:
	@echo "Available targets:"
	@echo "  generate     - Generate management API code from OpenAPI spec"
	@echo "  build        - Build the entire project"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean generated files"
	@echo "  install-tools - Install required tools for code generation"
	@echo "  all          - Run generate, build, and test"

# Install required tools for code generation
install-tools:
	@echo "Installing ogen for OpenAPI code generation..."
	go install github.com/ogen-go/ogen/cmd/ogen@latest

# Generate management API code from OpenAPI specification
generate:
	@echo "Generating management API code from OpenAPI spec..."
	cd kinde/management_api && go generate

# Build the entire project
build: generate
	@echo "Building the project..."
	go build ./...

# Run tests
test: generate
	@echo "Running tests..."
	go test ./...

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	cd kinde/management_api && rm -f oas_*.go

# Run all targets
all: generate build test

# Development workflow
dev: install-tools generate build test
	@echo "Development setup complete!"
