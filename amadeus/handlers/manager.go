package handlers

import (
	"net/http"
)

// RequestHandler provides an interface for all request handlers
type RequestHandler interface {
	Handle(w http.ResponseWriter, r *http.Request)
}
