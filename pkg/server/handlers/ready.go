package handlers

import (
	"net/http"
)

func Ready(f func() bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := http.StatusOK
		if !f() {
			code = http.StatusInternalServerError
		}
		w.WriteHeader(code)
	}
}
