#!/bin/bash

# Exit on error
set -e

# Load environment variables
if [ -f .env ]; then
  export $(cat .env | xargs)
fi

# Build the service
echo "Building sprint-service..."
cd "$(dirname "$0")"
go build -o ../../bin/sprint-service ./cmd/server

# Run the service
echo "Starting sprint-service..."
../../bin/sprint-service
