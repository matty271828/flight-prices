package search

import (
	"fmt"
	"strconv"
	"time"

	"github.com/matty271828/flight-prices/controller"
)

// Parameters for the search algorithms.
type Parameters struct {
	InitialTemperature float64
	CoolingRate        float64
}

type SimulatedAnnealing struct {
	ControllerManager controller.ControllerManager
	Params            *Parameters
}

// NewSimulatedAnnealing creates a new instance of SimulatedAnnealing with given parameters.
func NewSimulatedAnnealing(cm controller.ControllerManager, params *Parameters) *SimulatedAnnealing {
	return &SimulatedAnnealing{
		ControllerManager: cm,
		Params:            params,
	}
}

func (sa *SimulatedAnnealing) Run(origin, destination, departureDate string) Result {
	fmt.Println("Running simulated annealing...")

	price, err := sa.getFlightPrice(origin, destination, departureDate)
	if err != nil {
		fmt.Printf("error getting flight price: %s", err)
	}

	return Result{
		Date:  time.Now(),
		Price: price,
	}
}

func (sa *SimulatedAnnealing) getFlightPrice(origin, destination, departureDate string) (float64, error) {
	parsedOffers, err := sa.ControllerManager.FlightOffersSearch(origin, destination, departureDate, "2")
	if err != nil {
		return 0, fmt.Errorf("error conducting flight offers search%v \n", err)
	}

	// Ensure there are offers in the response
	if len(parsedOffers.Data) == 0 {
		return 0, fmt.Errorf("no flight offers available%v\n ", err)
	}

	// Convert GrandTotal (which is of string type) to float64
	price, err := strconv.ParseFloat(parsedOffers.Data[0].Price.GrandTotal, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing flight price: %v\n ", err)
	}

	return price, nil
}
