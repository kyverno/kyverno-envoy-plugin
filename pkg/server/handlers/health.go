package handlers

import (
	"net/http"
)

func Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Return a 200 OK status to indicate the server is healthy
		w.WriteHeader(http.StatusOK)
	}
}
