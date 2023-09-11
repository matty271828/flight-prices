package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/matty271828/flight-prices/controller"
)

type FISHandler struct {
	ControllerManager controller.ControllerManager
}

func NewFISHandler(cm controller.ControllerManager) *FISHandler {
	return &FISHandler{ControllerManager: cm}
}

func (h *FISHandler) HandleFlightInspirationSearch(w http.ResponseWriter, r *http.Request) {
	origin := r.URL.Query().Get("origin")

	if origin == "" {
		errorMsg := "Error: origin query parameter is required"
		log.Println(errorMsg)
		http.Error(w, errorMsg, http.StatusBadRequest)
		return
	}

	data, err := h.ControllerManager.FlightInspirationSearch(origin)
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
