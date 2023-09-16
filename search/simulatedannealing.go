package search

import (
	"fmt"
	"log"
	"math"
	"math/rand"
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

func (sa *SimulatedAnnealing) Run(origin, destination, startDate string, duration time.Duration) Result {
	fmt.Println("Running simulated annealing...")

	// Convert startDate to time.Time to calculate endDate
	layout := "2006-01-02"
	parsedStartDate, err := time.Parse(layout, startDate)
	if err != nil {
		fmt.Printf("could not parse startDate: %v", err)
		return Result{}
	}

	// Calculate endDate using the provided duration
	endDate := parsedStartDate.Add(duration)

	// Step 1: Start with an initial solution
	currentDate := randomDateInRange(parsedStartDate, endDate)
	bestDate := currentDate

	currentPrice, err := sa.getFlight(origin, destination, currentDate.Format("2006-01-02"))
	if err != nil {
		fmt.Printf("error getting initial flight price: %s", err)
		return Result{}
	}
	bestPrice := currentPrice
	log.Printf("Initial flight: Date: %s, Price: %f", currentDate, currentPrice)

	// Initial temperature and cooling rate from the Parameters
	T := sa.Params.InitialTemperature
	coolingRate := sa.Params.CoolingRate

	for T > 1 {
		// Step 3: Perturb the solution
		newDate := randomDateInRange(parsedStartDate, endDate)

		newPrice, err := sa.getFlight(origin, destination, newDate.Format("2006-01-02"))
		if err != nil {
			fmt.Printf("error getting perturbed flight price: %s", err)
			return Result{}
		}

		// Step 5: Acceptance criteria
		if newPrice != 0 && (newPrice < currentPrice || acceptNewSolution(currentPrice-newPrice, T)) {
			currentDate = newDate
			currentPrice = newPrice

			if currentPrice < bestPrice {
				bestPrice = currentPrice
				bestDate = currentDate
			}
		}
		log.Printf("Flight checked: Date: %s, Price: %f", currentDate, currentPrice)

		// Sleep for 2 seconds between each API call
		time.Sleep(time.Second * 2)

		// Step 6: Decrease the temperature
		T *= coolingRate
	}

	return Result{
		Date:  bestDate,
		Price: bestPrice,
	}
}

func (sa *SimulatedAnnealing) getFlight(origin, destination, departureDate string) (float64, error) {
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

// randomDateInRange returns a random date between start and end.
func randomDateInRange(start, end time.Time) time.Time {
	duration := end.Sub(start)
	randomDuration := time.Duration(rand.Int63n(int64(duration)))
	return start.Add(randomDuration)
}

// acceptNewSolution determines whether to accept a new solution based on the simulated annealing probability.
func acceptNewSolution(deltaCost float64, temperature float64) bool {
	if deltaCost < 0 {
		return true
	}
	return rand.Float64() < math.Exp(-deltaCost/temperature)
}
