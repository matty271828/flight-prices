package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/matty271828/flight-prices/amadeus"
	"github.com/matty271828/flight-prices/controller"
	"github.com/matty271828/flight-prices/server"
)

func main() {
	// Get the path of the currently running executable
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error determining executable path: %s\n", err)
	}
	basepath := filepath.Dir(execPath)

	amadeusClient, err := amadeus.NewAmadeusClient()
	if err != nil {
		err := fmt.Sprintf("Error getting amadeus client: %s\n", err)
		log.Println(err)
		return
	}

	c := controller.NewController(amadeusClient)
	s := server.NewServer(c)

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", s))

	uiPath := filepath.Join(basepath, "ui")
	log.Printf("Serving UI from: %s", uiPath) // Log the UI path for debugging
	fs := http.FileServer(http.Dir(uiPath))
	mux.Handle("/", fs)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", mux))
}
