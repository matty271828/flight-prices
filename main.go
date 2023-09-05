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

	// Setup environment variables and basepath
	basepath := setupEnv()

	// Create Amadeus Client
	amadeusClient, err := amadeus.NewAmadeusClient()
	if err != nil {
		errMsg := fmt.Sprintf("Error getting amadeus client: %s\n", err)
		// TODO: This is erroring on main - figure out why
		log.Println(errMsg)
	}

	c := controller.NewController(amadeusClient)

	// Setup servers
	setupServers(basepath, c, &wg)

	// Wait until all servers are done
	wg.Wait()
}

func setupEnv() string {
	// Attempt to load from .env file, if it exists
	_ = godotenv.Load()

	// Get the path of the currently running executable
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error determining executable path: %s\n", err)
	}
	return filepath.Dir(execPath)
}

func setupServers(basepath string, c controller.ControllerManager, wg *sync.WaitGroup) {
	// Init the first server
	s := server.NewServer(c)
	s.SetupRoutes()

	// Start UI Handler
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.Start(basepath, "ui", "/static/", "8080")
	}()

	// Init the dev server
	devS := server.NewServer(c)
	devS.SetupRoutes()

	// Start Dev UI Handler
	wg.Add(1)
	go func() {
		defer wg.Done()
		devS.Start(basepath, "ui-dev", "/devstatic/", "8091")
	}()
}
