package internal

import (
	"fmt"
	"net/http"

	"github.com/matty271828/flight-prices/search"
)

func TriggerSimulatedAnnealing(w http.ResponseWriter, r *http.Request) {
	sa := search.NewSimulatedAnnealing(&search.Parameters{
		InitialTemperature: 1000.0,
		CoolingRate:        0.95,
	})

	result := sa.Run()

	// Outputting to the console
	fmt.Printf("Optimal flight date: %v with price: %v\n", result.Date, result.Price)

	// Sending a response to the caller
	w.Write([]byte(fmt.Sprintf("Search complete. Optimal flight date: %s with price: €%.2f\n", result.Date.Format("2006-01-02"), result.Price)))

}
