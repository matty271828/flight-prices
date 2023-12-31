package controller

import (
	"github.com/matty271828/flight-prices/amadeus"
	"github.com/matty271828/flight-prices/amadeus/airportsearch"
	"github.com/matty271828/flight-prices/amadeus/flightinspiration"
	"github.com/matty271828/flight-prices/amadeus/flightoffers"
)

type Controller struct {
	AmadeusManager amadeus.AmadeusManager
}

func NewController(a amadeus.AmadeusManager) *Controller {
	return &Controller{AmadeusManager: a}
}

func (c *Controller) FlightInspirationSearch(origin string) (*flightinspiration.FISResponse, error) {
	return c.AmadeusManager.FlightInspirationSearch(origin)
}

func (c *Controller) FlightOffersSearch(origin, destination, departureDate, adults string) (*flightoffers.FOSResponse, error) {
	return c.AmadeusManager.FlightOffersSearch(origin, destination, departureDate, adults)
}

func (c *Controller) AirportSearch(name string) (*airportsearch.AirportSearchResponse, error) {
	return c.AmadeusManager.AirportSearch(name)
}
