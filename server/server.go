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
	server := &Server{
		ControllerManager: c,
		Router:            mux.NewRouter(),
	}
	return server
}

// SetupRoutes is called when we want to include
// the api on an initialised server instance.
func (s *Server) SetupRoutes() {
	apiRouter := s.Router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/get-destinations/", s.HandleGetDestinations).Methods("GET")
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
