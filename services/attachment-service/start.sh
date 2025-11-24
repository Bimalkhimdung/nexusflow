#!/bin/bash

# Exit on error
set -e

# Load environment variables
if [ -f .env ]; then
  export $(cat .env | xargs)
fi

# Build the service
echo "Building attachment-service..."
cd "$(dirname "$0")"
go build -o ../../bin/attachment-service ./cmd/server

# Run the service
echo "Starting attachment-service..."
../../bin/attachment-service
