package search

import (
	"errors"
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

func (sa *SimulatedAnnealing) Run(origin, destination, departureDate string) (Result, error) {
	fmt.Println("Running simulated annealing...")

	price, _ := sa.getFlightPrice(origin, destination, departureDate)

	return Result{
		Date:  time.Now(),
		Price: price,
	}, nil
}

func (sa *SimulatedAnnealing) getFlightPrice(origin, destination, date string) (float64, error) {
	parsedOffers, err := sa.ControllerManager.FlightOffersSearch(origin, destination, "1", date)
	if err != nil {
		return 0, err
	}

	// Ensure there are offers in the response
	if len(parsedOffers.Data) == 0 {
		return 0, errors.New("no flight offers available")
	}

	// Convert GrandTotal (which is of string type) to float64
	price, err := strconv.ParseFloat(parsedOffers.Data[0].Price.GrandTotal, 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing flight price: %v", err)
	}

	return price, nil
}
