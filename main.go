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

	// Attempt to load from .env file, if it exists
	_ = godotenv.Load() // We ignore the error since it's optional

	// Get the path of the currently running executable
	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error determining executable path: %s\n", err)
	}
	basepath := filepath.Dir(execPath)

	fmt.Println(basepath)

	cfg := amadeus.Config{
		ClientId:     os.Getenv("AMADEUS_API_KEY"),
		ClientSecret: os.Getenv("AMADEUS_API_SECRET"),
	}

	amadeusClient, err := amadeus.NewAmadeusClient(cfg)
	if err != nil {
		err := fmt.Sprintf("Error getting amadeus client: %s\n", err)
		log.Println(err)
		// don't return error as we want the app to still run if
		// cannot connect to amadeus.
	}

	c := controller.NewController(amadeusClient)

	s := server.NewServer(c)
	devS := server.NewServer(c)
	devS.SetupRoutes() // Only setup API routes for the 8091 server

	// Start UI Handler
	wg.Add(1)
	go func() {
		defer wg.Done()
		setupServer(basepath, "ui", "/static/", "8080", s)
	}()

	// Start Dev UI Handler
	wg.Add(1)
	go func() {
		defer wg.Done()
		setupServer(basepath, "ui-dev", "/devstatic/", "8091", devS)
	}()

	// Wait until all servers are done
	wg.Wait()
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
