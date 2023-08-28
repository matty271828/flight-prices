package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/matty271828/flight-prices/controller"
)

type Server struct {
	ControllerManager controller.ControllerManager
}

func NewServer(c controller.ControllerManager) *Server {
	return &Server{ControllerManager: c}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("Request received for path: %s\n", r.URL.Path)
	log.Println(msg)

	switch r.URL.Path {
	case "/get-destinations/":
		s.handleGetDestinations(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) handleGetDestinations(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	origin := query.Get("origin")

	if origin == "" {
		http.Error(w, "origin query parameter is required", http.StatusBadRequest)
		log.Println("Error: origin query parameter is required")
		return
	}

	data, err := s.ControllerManager.GetFlightInfo(origin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		err := fmt.Sprintf("Error getting flight info: %v\n", err)
		log.Println(err)
		return
	}

	// Marshalling data with indentation
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		err := fmt.Sprintf("Error marshalling data: %v\n", err)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		err := fmt.Sprintf("Error writing response: %v\n", err)
		log.Println(err)
	}
}
