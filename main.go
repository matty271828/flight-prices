package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
	"github.com/matty271828/flight-prices/amadeus"
	"github.com/matty271828/flight-prices/controller"
	"github.com/matty271828/flight-prices/server"
)

func main() {
	var wg sync.WaitGroup

	basepath, err := setupEnv()
	if err != nil {
		log.Fatalf("Failed to set up environment: %v", err)
	}

	// Create Amadeus Client
	amadeusClient, err := amadeus.NewAmadeusClient()
	if err != nil {
		errMsg := fmt.Sprintf("Error getting amadeus client: %s\n", err)
		// TODO: This is erroring on main - figure out why
		log.Fatalf(errMsg)
	}

	c := controller.NewController(amadeusClient)

	// Initialize main server
	if _, err = server.NewServer(c, basepath, "ui", "/static/", "8080", &wg); err != nil {
		log.Fatalf("Failed to initialize main server: %v", err)
	}

	// Initialize dev server
	if _, err = server.NewServer(c, basepath, "ui-dev", "/devstatic/", "8091", &wg); err != nil {
		log.Fatalf("Failed to initialize dev server: %v", err)
	}

	// Wait until all servers are done
	wg.Wait()
}

func setupEnv() (string, error) {
	// Attempt to load from .env file, if it exists
	_ = godotenv.Load()

	// Get the path of the currently running executable
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error determining executable path: %s\n", err)
		return "", fmt.Errorf("Error determining executable path: %w", err)
	}
	return filepath.Dir(execPath), nil
}
