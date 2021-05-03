package internal

import (
	"net/http"
	"time"
)

func WriteResponse(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(message))
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func MakeServerFromMux(mux http.Handler) *http.Server {
	// set timeouts so that a slow or malicious client doesn't
	// hold resources forever
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}
