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

	"github.com/matty271828/flight-prices/amadeus"
	"github.com/matty271828/flight-prices/controller"
	"github.com/matty271828/flight-prices/server"
)

func main() {
	var wg sync.WaitGroup

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

	// Start UI Handler
	wg.Add(1)
	go func() {
		defer wg.Done()
		routes := []string{"/api"}
		setupServer(basepath, "ui", "8080", routes, s)
	}()

	// Start Dev UI Handler
	wg.Add(1)
	go func() {
		defer wg.Done()
		routes := []string{}
		setupServer(basepath, "ui-dev", "8091", routes, s)
	}()

	// Wait until all servers are done
	wg.Wait()
}

// setupServer is used to setup a server on a requested port with a supplied ui
func setupServer(basepath, dir, port string, routes []string, s *server.Server) {
	mux := http.NewServeMux()
	if routes != nil {
		for _, route := range routes {
			mux.Handle(route+"/", http.StripPrefix(route, s))
		}
	}

	// Serve index.html dynamically, this is for cachebusting
	mux.HandleFunc("/", serveIndexFromDir(dir))

	dirPath := filepath.Join(basepath, dir)
	log.Printf("Serving UI from: %s", dirPath) // Log the dirPath for debugging

	// Serve other static files
	fs := http.FileServer(http.Dir(dirPath))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
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
