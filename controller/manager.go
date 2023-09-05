package controller

import "github.com/matty271828/flight-prices/amadeus"

type ControllerManager interface {
	FlightInspirationSearch(origin string) (*amadeus.ApiResponse, error)
}
