#!/bin/bash

# Exit on error
set -e

# Load environment variables
if [ -f .env ]; then
  export $(cat .env | xargs)
fi

# Build the service
echo "Building board-service..."
cd "$(dirname "$0")"
go build -o ../../bin/board-service ./cmd/server

# Run the service
echo "Starting board-service..."
../../bin/board-service
