package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Middleware mux.MiddlewareFunc

// Chain is used to add multiple middlewares to a router
func Chain(middlewares ...mux.MiddlewareFunc) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}
