package controller

import "github.com/matty271828/flight-prices/amadeus"

type ControllerManager interface {
	GetFlightInfo(origin string) (*amadeus.ApiResponse, error)
}
