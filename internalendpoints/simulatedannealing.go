package internal

import (
	"fmt"
	"net/http"
	"time"

	"github.com/matty271828/flight-prices/controller"
	"github.com/matty271828/flight-prices/search"
)

type SAHandler struct {
	ControllerManager controller.ControllerManager
}

func NewSAHandler(c controller.ControllerManager) *SAHandler {
	return &SAHandler{ControllerManager: c}
}

func (h *SAHandler) Handle(w http.ResponseWriter, r *http.Request) {
	origin := r.URL.Query().Get("origin")
	destination := r.URL.Query().Get("destination")
	dateStr := r.URL.Query().Get("date")

	// Validate the parameters
	if origin == "" || destination == "" || dateStr == "" {
		w.Write([]byte("Please provide valid origin, destination, and date."))
		return
	}

	sa := search.NewSimulatedAnnealing(h.ControllerManager, &search.Parameters{
		InitialTemperature: 1000.0,
		CoolingRate:        0.95,
	})

	// specifiy period to search over
	period := 1 * 30 * time.Hour * 24
	result := sa.Run(origin, destination, dateStr, period)

	// Outputting to the console
	fmt.Printf("Optimal flight date: %v with price: %v\n", result.Date, result.Price)

	// Sending a response to the caller
	w.Write([]byte(fmt.Sprintf("Search complete. Optimal flight date: %s with price: €%.2f\n", result.Date.Format("2006-01-02"), result.Price)))
}
