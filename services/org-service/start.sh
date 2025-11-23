#!/bin/bash

# Set environment variables
export ORG_SERVICE_SERVER_GRPC_PORT=50052
export ORG_SERVICE_KAFKA_BROKERS=localhost:19092

# Run the service
cd "$(dirname "$0")"
../../bin/org-service
