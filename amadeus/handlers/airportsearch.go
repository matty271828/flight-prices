package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/matty271828/flight-prices/controller"
)

type AirportSearchHandler struct {
	ControllerManager controller.ControllerManager
}

func NewAirportSearchHandler(cm controller.ControllerManager) *AirportSearchHandler {
	return &AirportSearchHandler{ControllerManager: cm}
}

// Ensure AirportSearchHandler implements the RequestHandler interface.
var _ RequestHandler = (*AirportSearchHandler)(nil)

func (h *AirportSearchHandler) Handle(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")

	if keyword == "" {
		errorMsg := "Error: keyword query parameter is required"
		log.Println(errorMsg)
		http.Error(w, errorMsg, http.StatusBadRequest)
		return
	}

	data, err := h.ControllerManager.AirportSearch(keyword)
	if err != nil {
		errorMsg := fmt.Sprintf("Error getting airport info for keyword %s: %v", keyword, err)
		log.Println(errorMsg)
		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		errorMsg := fmt.Sprintf("Error marshalling data for origin %s: %v", keyword, err)
		log.Println(errorMsg)
		http.Error(w, errorMsg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonData)
	if err != nil {
		errorMsg := fmt.Sprintf("Error writing response for origin %s: %v", keyword, err)
		log.Println(errorMsg)
		return
	}
}
