#!/bin/bash

# Exit on error
set -e

# Load environment variables
if [ -f .env ]; then
  export $(cat .env | xargs)
fi

# Build the service
echo "Building notification-service..."
cd "$(dirname "$0")"
go build -o ../../bin/notification-service ./cmd/server

# Run the service
echo "Starting notification-service..."
../../bin/notification-service
