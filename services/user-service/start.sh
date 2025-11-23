#!/bin/bash

# Set environment variables
export USER_SERVICE_SERVER_GRPC_PORT=50051
export USER_SERVICE_KAFKA_BROKERS=localhost:19092

# Run the service
cd "$(dirname "$0")"
../../bin/user-service
