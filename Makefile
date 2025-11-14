.PHONY: build clean install test fmt vet

# Binary name
BINARY=linesense

# Build the binary
build:
	go build -o $(BINARY) ./cmd/linesense

# Clean build artifacts
clean:
	rm -f $(BINARY)
	go clean

# Install to $GOPATH/bin or $HOME/go/bin
install:
	go install ./cmd/linesense

# Run tests (when implemented)
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run all checks
check: fmt vet

# Build for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build -o $(BINARY)-linux-amd64 ./cmd/linesense
	GOOS=darwin GOARCH=amd64 go build -o $(BINARY)-darwin-amd64 ./cmd/linesense
	GOOS=darwin GOARCH=arm64 go build -o $(BINARY)-darwin-arm64 ./cmd/linesense
	GOOS=windows GOARCH=amd64 go build -o $(BINARY)-windows-amd64.exe ./cmd/linesense

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build the binary"
	@echo "  clean      - Remove build artifacts"
	@echo "  install    - Install to GOPATH/bin"
	@echo "  test       - Run tests"
	@echo "  fmt        - Format code"
	@echo "  vet        - Run go vet"
	@echo "  check      - Run fmt and vet"
	@echo "  build-all  - Build for multiple platforms"
	@echo "  help       - Show this help message"
