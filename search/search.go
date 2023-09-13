package search

import "time"

// Strategy represents a search optimization algorithm.
type Strategy interface {
	Run() Result
}

// Result represents the outcome of a search optimization.
type Result struct {
	Date  time.Time
	Price float64
}
