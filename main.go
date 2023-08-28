package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/matty271828/flight-prices/amadeus"
	"github.com/matty271828/flight-prices/controller"
	"github.com/matty271828/flight-prices/server"
)

func main() {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	amadeusClient, err := amadeus.NewAmadeusClient()
	if err != nil {
		err := fmt.Sprintf("Error getting amadeus client: %s\n", err)
		log.Println(err)
		return
	}

	c := controller.NewController(amadeusClient)
	s := server.NewServer(c)

	mux := http.NewServeMux()
	mux.Handle("/api/", http.StripPrefix("/api", s))

	fs := http.FileServer(http.Dir(filepath.Join(basepath, "ui")))
	mux.Handle("/", fs)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", mux))
}
