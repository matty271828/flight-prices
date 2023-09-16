package internal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/matty271828/flight-prices/controller"
	"github.com/matty271828/flight-prices/search"
)

type SSHandler struct {
	ControllerManager controller.ControllerManager
}

func NewSSHandler(c controller.ControllerManager) *SSHandler {
	return &SSHandler{ControllerManager: c}
}

func (h *SSHandler) Handle(w http.ResponseWriter, r *http.Request) {
	origin := r.URL.Query().Get("origin")
	destination := r.URL.Query().Get("destination")

	// Validate the parameters
	if origin == "" || destination == "" {
		w.Write([]byte("Please provide valid origin and destination"))
		return
	}

	// Set the interval here, number of days to jump between searches
	ss := search.NewStructuredSearch(h.ControllerManager, 18)

	// specifiy period to search over
	period := 6 * 30 * time.Hour * 24
	// TODO: Retrieve current date and pass in as startDate for run
	result := ss.Run(origin, destination, "2023-10-01", period)

	// Outputting to the console
	fmt.Printf("%s to %s : Cheapest flight found : %v with price: %v\n", origin, destination, result.Date, result.Price)

	// Sending a response to the caller
	w.Write([]byte(fmt.Sprintf("Structured search complete: Cheapest flight found: %s with price: €%.2f\n", result.Date.Format("2006-01-02"), result.Price)))
}
