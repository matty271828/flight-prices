package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matty271828/flight-prices/controller"
)

type Server struct {
	ControllerManager controller.ControllerManager
	Router            *mux.Router
}

func NewServer(c controller.ControllerManager) *Server {
	server := &Server{ControllerManager: c}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.Router = mux.NewRouter()
	s.Router.HandleFunc("/get-destinations/", s.HandleGetDestinations).Methods("GET")
}

func (s *Server) HandleGetDestinations(w http.ResponseWriter, r *http.Request) {
	origin := r.URL.Query().Get("origin")

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
