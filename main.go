package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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
		log.Println(errMsg)
	}

	c := controller.NewController(amadeusClient)

	err = setupServers(basepath, c, &wg)
	if err != nil {
		log.Fatalf("Failed to set up servers: %v", err)
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

func setupServers(basepath string, c controller.ControllerManager, wg *sync.WaitGroup) error {
	// Use error channels to capture errors from goroutines
	errCh := make(chan error, 2) // buffered channel to avoid potential deadlocks

	// Init the first server
	s := server.NewServer(c)
	s.SetupRoutes()

	// Start UI Handler
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.Start(basepath, "ui", "/static/", "8080"); err != nil {
			errCh <- fmt.Errorf("Failed to start main server: %w", err)
			return
		}
	}()

	// Init the dev server
	devS := server.NewServer(c)
	devS.SetupRoutes()

	// Start Dev UI Handler
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := devS.Start(basepath, "ui-dev", "/devstatic/", "8091"); err != nil {
			errCh <- fmt.Errorf("Failed to start dev server: %w", err)
			return
		}
	}()

	// Wait for goroutines to finish and close the error channel
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Collect any errors from the channel
	var errs []string
	for err := range errCh {
		errs = append(errs, err.Error())
	}
	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, " | "))
	}
	return nil
}
