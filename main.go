package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

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

func setupServers(basepath string, c *controller.Controller, wg *sync.WaitGroup) {
	s := server.NewServer(c)
	devS := server.NewServer(c)
	devS.SetupRoutes()

	// Start UI Handler
	wg.Add(1)
	go startServer(basepath, "ui", "/static/", "8080", s, wg)

	// Start Dev UI Handler
	wg.Add(1)
	go startServer(basepath, "ui-dev", "/devstatic/", "8091", devS, wg)
}

func startServer(basepath, uiType, route, port string, s *server.Server, wg *sync.WaitGroup) {
	defer wg.Done()
	setupServer(basepath, uiType, route, port, s)
}

// setupServer is used to setup a server on a requested port with a supplied ui
func setupServer(basepath, dir, staticRoute, port string, s *server.Server) {
	r := s.Router

	// Serve the index.html dynamically for cachebusting
	r.HandleFunc("/", serveIndexFromDir(dir))

	// Serve other static files
	dirPath := filepath.Join(basepath, dir)
	fs := http.FileServer(http.Dir(dirPath))
	r.PathPrefix(staticRoute).Handler(http.StripPrefix(staticRoute, fs))

	log.Printf("HTTP Server initialized on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// serveIndexFromDir is used to serve the index file of a ui directory
func serveIndexFromDir(dirPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		modifiedContent, err := generateHTML(filepath.Join(dirPath, "index.html"))
		if err != nil {
			http.Error(w, "Failed to generate content", http.StatusInternalServerError)
			return
		}
		w.Write([]byte(modifiedContent))
	}
}

func generateHTML(filepath string) (string, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	// Add timestamps for cachebusting
	modifiedContent := strings.ReplaceAll(string(content), "{{cssTimestamp}}", timestamp)
	modifiedContent = strings.ReplaceAll(modifiedContent, "{{jsTimestamp}}", timestamp)

	return modifiedContent, nil
}
