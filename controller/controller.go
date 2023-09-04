package controller

import (
	"github.com/matty271828/flight-prices/amadeus"
)

type Controller struct {
	AmadeusManager amadeus.AmadeusManager
}

func NewController(a amadeus.AmadeusManager) *Controller {
	return &Controller{AmadeusManager: a}
}

func (c *Controller) FlightOffersSearch(origin string) (*amadeus.ApiResponse, error) {
	return c.AmadeusManager.FlightOffersSearch(origin)
}
