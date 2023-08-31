package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	// Serve index.html dynamically, this is for cachebusting
	mux.HandleFunc("/", indexHandler)

	uiPath := filepath.Join(basepath, "ui")
	log.Printf("Serving UI from: %s", uiPath) // Log the UI path for debugging

	// Serve other static files
	fs := http.FileServer(http.Dir(uiPath))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	modifiedContent, err := generateHTML()
	if err != nil {
		http.Error(w, "Failed to generate content", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(modifiedContent))
}

func generateHTML() (string, error) {
	content, err := ioutil.ReadFile("ui/index.html")
	if err != nil {
		return "", err
	}

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	// Add timestamps for cachebusting
	modifiedContent := strings.ReplaceAll(string(content), "{{cssTimestamp}}", timestamp)
	modifiedContent = strings.ReplaceAll(modifiedContent, "{{jsTimestamp}}", timestamp)

	return modifiedContent, nil
}
