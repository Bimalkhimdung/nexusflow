#!/bin/bash

# Exit on error
set -e

# Load environment variables
if [ -f .env ]; then
  export $(cat .env | xargs)
fi

# Build the service
echo "Building issue-service..."
cd "$(dirname "$0")"
go build -o ../../bin/issue-service ./cmd/server

# Run the service
echo "Starting issue-service..."
../../bin/issue-service
