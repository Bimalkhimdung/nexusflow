#!/bin/bash

# Exit on error
set -e

# Load environment variables
if [ -f .env ]; then
  export $(cat .env | xargs)
fi

# Build the service
echo "Building search-service..."
cd "$(dirname "$0")"
go build -o ../../bin/search-service ./cmd/server

# Run the service
echo "Starting search-service..."
../../bin/search-service
