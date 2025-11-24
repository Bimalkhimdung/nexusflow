#!/bin/bash

# Script to generate Go code from protobuf definitions

set -e

echo "Generating protobuf code..."

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed"
    echo "Install with: brew install protobuf (macOS) or apt-get install protobuf-compiler (Linux)"
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Error: protoc-gen-go is not installed"
    echo "Install with: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    exit 1
fi

# Check if protoc-gen-grpc-gateway is installed
if ! command -v protoc-gen-grpc-gateway &> /dev/null; then
    echo "Installing protoc-gen-grpc-gateway..."
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
    go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
fi

# Create output directory
rm -rf pkg/proto/proto
mkdir -p pkg/proto

# Find all .proto files (excluding google/api)
PROTO_FILES=$(find proto -name "*.proto" | grep -v "proto/google")
PKG_PROTO_FILES=$(find pkg/proto -name "*.proto")

# Generate Go code for each proto file
for proto_file in $PROTO_FILES $PKG_PROTO_FILES; do
    echo -e "${BLUE}Generating code for $proto_file${NC}"
    
    protoc \
        --go_out=. \
        --go_opt=module=github.com/nexusflow/nexusflow \
        --go-grpc_out=. \
        --go-grpc_opt=module=github.com/nexusflow/nexusflow \
        --grpc-gateway_out=. \
        --grpc-gateway_opt=module=github.com/nexusflow/nexusflow \
        --proto_path=. \
        --proto_path=proto \
        "$proto_file"
done

echo -e "${GREEN}✓ Protobuf code generation complete!${NC}"
echo "Generated files are in pkg/proto/"
