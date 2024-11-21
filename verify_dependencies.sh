#!/bin/bash
echo "Verifying dependencies..."
go mod tidy
go mod verify
if [ $? -eq 0 ]; then
    echo "Dependencies verified successfully!"
else
    echo "Dependency verification failed. Please fix the issues."
    exit 1
fi
