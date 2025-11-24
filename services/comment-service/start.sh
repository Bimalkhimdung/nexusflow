#!/bin/bash

# Exit on error
set -e

# Load environment variables
if [ -f .env ]; then
  export $(cat .env | xargs)
fi

# Build the service
echo "Building comment-service..."
cd "$(dirname "$0")"
go build -o ../../bin/comment-service ./cmd/server

# Run the service
echo "Starting comment-service..."
../../bin/comment-service
