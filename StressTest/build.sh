#!/bin/bash

echo "Building stresstest application..."

# Build the Go application
go build -o stresstest .

if [ $? -eq 0 ]; then
    echo "Build successful! Binary created: stresstest"
    echo ""
    echo "Usage examples:"
    echo "  ./stresstest --url=http://google.com --requests=100 --concurrency=10"
    echo "  ./stresstest --url=https://httpbin.org/get --requests=50 --concurrency=5"
else
    echo "Build failed!"
    exit 1
fi
