#!/bin/bash

# Exit on error
set -e

# Load environment variables
if [ -f .env ]; then
  export $(cat .env | xargs)
fi

# Build the service
echo "Building workflow-service..."
cd "$(dirname "$0")"
go build -o ../../bin/workflow-service ./cmd/server

# Run the service
echo "Starting workflow-service..."
../../bin/workflow-service
