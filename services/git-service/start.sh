#!/bin/bash

# Exit on error
set -e

# Load environment variables
if [ -f .env ]; then
  export $(cat .env | xargs)
fi

# Build the service
echo "Building git-service..."
cd "$(dirname "$0")"
go build -o ../../bin/git-service ./cmd/server

# Run the service
echo "Starting git-service..."
../../bin/git-service
