package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/matty271828/flight-prices/controller"
)

type Server struct {
	ControllerManager controller.ControllerManager
	Router            *mux.Router
	UIBasepath        string
	UIType            string
	Route             string
	Port              string
}

func NewServer(c controller.ControllerManager) *Server {
	server := &Server{
		ControllerManager: c,
		Router:            mux.NewRouter(),
	}
	return server
}

// SetupRoutes is called when we want to include the API on an initialised server instance.
func (s *Server) SetupRoutes() {
	apiRouter := s.Router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/get-destinations/", s.HandleGetDestinations).Methods("GET")
}

func (s *Server) Start(basepath, uiType, route, port string) error {
	s.UIBasepath = basepath
	s.UIType = uiType
	s.Route = route
	s.Port = port

	if err := s.setupUIRoutes(); err != nil {
		return fmt.Errorf("Error setting up UI routes: %w", err)
	}

	if err := s.setupServer(); err != nil {
		return fmt.Errorf("Error setting up server: %w", err)
	}

	return nil
}
func (s *Server) setupUIRoutes() error {
	// Serve the index.html dynamically for cachebusting
	s.Router.HandleFunc("/", s.serveIndex())

	// Serve other static files
	dirPath := filepath.Join(s.UIBasepath, s.UIType)

	// Check if dirPath exists or is accessible.
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return fmt.Errorf("Directory %s does not exist", dirPath)
	}

	fs := http.FileServer(http.Dir(dirPath))
	s.Router.PathPrefix(s.Route).Handler(http.StripPrefix(s.Route, fs))
	return nil
}

func (s *Server) setupServer() error {
	log.Printf("HTTP Server initialized on port %s", s.Port)
	err := http.ListenAndServe(":"+s.Port, s.Router)
	if err != nil {
		return fmt.Errorf("Failed to start server on port %s: %w", s.Port, err)
	}
	return nil
}

func (s *Server) serveIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		modifiedContent, err := generateHTML(filepath.Join(s.UIType, "index.html"))
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

func (s *Server) HandleGetDestinations(w http.ResponseWriter, r *http.Request) {
	origin := r.URL.Query().Get("origin")

	if origin == "" {
		errorMsg := "Error: origin query parameter is required"
		log.Println(errorMsg)
		http.Error(w, errorMsg, http.StatusBadRequest)
		return
	}

	data, err := s.ControllerManager.FlightOffersSearch(origin)
	if err != nil {
		errorMsg := fmt.Sprintf("Error getting flight info for origin %s: %v", origin, err)
		log.Println(errorMsg)
		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		errorMsg := fmt.Sprintf("Error marshalling data for origin %s: %v", origin, err)
		log.Println(errorMsg)
		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		errorMsg := fmt.Sprintf("Error writing response for origin %s: %v", origin, err)
		log.Println(errorMsg)
		return
	}
}
