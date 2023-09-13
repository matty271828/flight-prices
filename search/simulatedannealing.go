package search

import (
	"fmt"
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

func (sa *SimulatedAnnealing) Run(origin, destination string, departureDate time.Time) Result {
	fmt.Println("Running simulated annealing...")
	// Return a mock result
	return Result{
		Date:  time.Now(),
		Price: 100.50,
	}
}
