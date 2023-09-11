package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/matty271828/flight-prices/controller"
)

type FOSHandler struct {
	ControllerManager controller.ControllerManager
}

func NewFOSHandler(cm controller.ControllerManager) *FOSHandler {
	return &FOSHandler{ControllerManager: cm}
}

func (h *FOSHandler) HandleFlightOffersSearch(w http.ResponseWriter, r *http.Request) {
	requiredParams := []string{"origin", "destination", "departureDate", "adults"}
	params := make(map[string]string)

	// Loop through the required parameters and check if they are present
	for _, param := range requiredParams {
		value := r.URL.Query().Get(param)
		if value == "" {
			errorMsg := fmt.Sprintf("Error: %s query parameter is required", param)
			log.Println(errorMsg)
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}
		params[param] = value
	}

	data, err := h.ControllerManager.FlightOffersSearch(
		params["origin"],
		params["destination"],
		params["departureDate"],
		params["adults"],
	)
	if err != nil {
		errorMsg := fmt.Sprintf("Error getting flight offers for origin %s: %v", params["origin"], err)
		log.Println(errorMsg)
		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		errorMsg := fmt.Sprintf("Error marshalling data for origin %s: %v", params["origin"], err)
		log.Println(errorMsg)
		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		errorMsg := fmt.Sprintf("Error writing response for origin %s: %v", params["origin"], err)
		log.Println(errorMsg)
		return
	}
}
