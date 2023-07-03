package server

import (
	"encoding/json"
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
	switch r.URL.Path {
	case "/":
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
		return
	}

	data, err := s.ControllerManager.GetFlightInfo(origin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Marshalling data with indentation
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
