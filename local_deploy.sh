#!/bin/bash

# Step 1: Build the binary
echo "Building the application..."
go build -o flight-prices-binary

# Check for successful build
if [ $? -ne 0 ]; then
    echo "Build failed. Exiting."
    exit 1
fi

# Step 2: Run the binary
echo "Starting the application..."
./flight-prices-binary & echo $! > app.pid

# Allow some time for the server to start
sleep 2

# Step 3: Open a web browser
echo "Opening web browser..."
case "$OSTYPE" in
  darwin*) open "http://localhost:8080" ;; 
  linux*) xdg-open "http://localhost:8080" ;;
  *) echo "Unknown OS type. Please open a web browser manually at http://localhost:8080" ;;
esac

exit 0
