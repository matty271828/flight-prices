#!/bin/bash

# Step 1: Load environment variables from .env file
source .env

# Step 2: Build the binary
echo "Building the application..."
go build -o flight-prices-binary

# Check for successful build
if [ $? -ne 0 ]; then
    echo "Build failed. Exiting."
    exit 1
fi

# Check if a process is already running
if [ -f app.pid ]; then
    pid=$(cat app.pid)
    if ps -p $pid > /dev/null; then
        echo "Stopping the existing application (PID: $pid)..."
        kill $pid
        sleep 1
    fi
    rm app.pid
fi

# Step 3: Run the binary
echo "Starting the application..."
./flight-prices-binary & echo $! > app.pid

# Allow some time for the server to start
sleep 2

# Step 4: Open a web browser
echo "Opening web browser..."
case "$OSTYPE" in
  darwin*) open "http://localhost:8080" ;; 
  linux*) xdg-open "http://localhost:8080" ;;
  *) echo "Unknown OS type. Please open a web browser manually at http://localhost:8080" ;;
esac

exit 0
