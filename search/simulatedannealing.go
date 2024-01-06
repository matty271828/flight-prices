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

	parsedStartDate, endDate, err := sa.parseAndComputeDates(startDate, duration)
	if err != nil {
		fmt.Printf("could not parse startDate: %v", err)
		return Result{}
	}

	currentDate, currentPrice, err := sa.initializeSolution(origin, destination, parsedStartDate, endDate)
	if err != nil {
		fmt.Printf("error getting initial flight price: %s", err)
		return Result{}
	}

	bestDate, bestPrice := currentDate, currentPrice

	T := sa.Params.InitialTemperature
	coolingRate := sa.Params.CoolingRate

	maxIterations := 5
	iterationCount := 0
	totalAPICost := 0.0
	apiCallCost := 0.025

	for T > 1 && iterationCount < maxIterations {
		currentDate, currentPrice = sa.anneal(origin, destination, parsedStartDate, endDate, currentPrice, currentDate, T)

		if currentPrice < bestPrice {
			bestPrice = currentPrice
			bestDate = currentDate
		}

		iterationCount++
		totalAPICost = float64(iterationCount) * apiCallCost
		log.Printf("Flight checked: Date: %s, Price: %f", currentDate, currentPrice)
		fmt.Printf("Iteration: %d, Total API Cost: €%.2f\n", iterationCount, totalAPICost)

		time.Sleep(time.Second * 2)
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

// Helper function to parse date and calculate end date
func (sa *SimulatedAnnealing) parseAndComputeDates(startDate string, duration time.Duration) (time.Time, time.Time, error) {
	layout := "2006-01-02"
	parsedStartDate, err := time.Parse(layout, startDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	endDate := parsedStartDate.Add(duration)
	return parsedStartDate, endDate, nil
}

// Helper function to initialize the solution
func (sa *SimulatedAnnealing) initializeSolution(origin, destination string, startDate, endDate time.Time) (time.Time, float64, error) {
	currentDate := randomDateInRange(startDate, endDate)
	currentPrice, err := sa.getFlight(origin, destination, currentDate.Format("2006-01-02"))
	if err != nil {
		return currentDate, 0, err
	}
	return currentDate, currentPrice, nil
}

// Helper function for the annealing process
func (sa *SimulatedAnnealing) anneal(
	origin string,
	destination string,
	startDate, endDate time.Time,
	currentPrice float64,
	currentDate time.Time,
	temperature float64,
) (time.Time, float64) {
	newDate := randomDateInRange(startDate, endDate)
	newPrice, err := sa.getFlight(origin, destination, newDate.Format("2006-01-02"))
	if err != nil {
		return currentDate, currentPrice
	}

	if newPrice != 0 && (newPrice < currentPrice || acceptNewSolution(currentPrice-newPrice, temperature)) {
		return newDate, newPrice
	}
	return currentDate, currentPrice
}
