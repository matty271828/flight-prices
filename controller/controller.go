package controller

import (
	"github.com/matty271828/flight-prices/amadeus"
)

type Controller struct {
	AmadeusClient *amadeus.AmadeusClient
}

func NewController(a *amadeus.AmadeusClient) *Controller {
	return &Controller{AmadeusClient: a}
}

func (c *Controller) GetFlightInfo(origin string) (*amadeus.ApiResponse, error) {
	return c.AmadeusClient.GetFlightInfo(origin)
}
