
# Set the target binary name and directory
BINARY_NAME=./bin
MYGOFLAGS=-ldflags "-s -w" # Add your custom build flags here

# Main build target
all: fmt build

# Format the code using gofmt
fmt:
	go mod tidy
	go fmt ./...

# Build the binary
build:
	mkdir -p ./bin
	go mod tidy
	go build $(MYGOFLAGS) -o $(BINARY_NAME) ./...

trivy:
	trivy --skip-update fs . --scanners vuln
