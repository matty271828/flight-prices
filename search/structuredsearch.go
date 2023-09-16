package search

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/matty271828/flight-prices/controller"
)

type StructuredSearch struct {
	ControllerManager controller.ControllerManager
	IntervalDays      int
}

func NewStructuredSearch(cm controller.ControllerManager, interval int) *StructuredSearch {
	return &StructuredSearch{
		ControllerManager: cm,
		IntervalDays:      interval,
	}
}

func (ss *StructuredSearch) Run(origin, destination, startDate string, duration time.Duration) Result {
	fmt.Println("Running structured search...")

	// Parse the start date and compute the end date
	parsedStartDate, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		fmt.Printf("could not parse startDate: %v", err)
		return Result{}
	}
	endDate := parsedStartDate.Add(duration)

	// Initializing best price with a very high value and zero date
	bestPrice := math.MaxFloat64
	var bestDate time.Time

	// Setting up the search intervals
	intervalDays := ss.IntervalDays
	apiCallCost := 0.025
	totalAPICost := 0.0

	// Define the maximum number of API calls
	maxAPICalls := 10
	apiCallCount := 0

	currentDate := parsedStartDate
	for currentDate.Before(endDate) && apiCallCount < maxAPICalls {
		price, err := ss.getFlight(origin, destination, currentDate.Format("2006-01-02"))
		if err != nil {
			fmt.Printf("Error fetching flight for date %s: %v\n", currentDate, err)
			continue
		}

		apiCallCount++
		totalAPICost = float64(apiCallCount) * apiCallCost

		// Update the best price and date if needed
		if price != 0 && price < bestPrice {
			bestPrice = price
			bestDate = currentDate
		}

		log.Printf("Flight checked: Date: %s, Price: €%.2f", currentDate, price)
		fmt.Printf("Iteration: %d, Total API Cost: €%.2f\n", apiCallCount, totalAPICost)

		// Move to the next interval date
		currentDate = currentDate.AddDate(0, 0, intervalDays)
	}

	return Result{
		Date:  bestDate,
		Price: bestPrice,
	}
}

func (ss *StructuredSearch) getFlight(origin, destination, departureDate string) (float64, error) {
	parsedOffers, err := ss.ControllerManager.FlightOffersSearch(origin, destination, departureDate, "2")
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
