#!/bin/bash

# Exit on error
set -e

# Load environment variables
if [ -f .env ]; then
  export $(cat .env | xargs)
fi

# Build the service
echo "Building project-service..."
cd "$(dirname "$0")"
go build -o ../../bin/project-service ./cmd/server

# Run the service
echo "Starting project-service..."
../../bin/project-service
