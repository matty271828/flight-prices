package server

import (
	"encoding/json"
	"net/http"

	"github.com/matty271828/flight-prices/controller"
)

type Server struct {
	Controller *controller.Controller
}

func NewServer(c *controller.Controller) *Server {
	return &Server{Controller: c}
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

	destinations, err := s.Controller.GetDestinations(origin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(destinations)
}
