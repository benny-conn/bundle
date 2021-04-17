package internal

import (
	"net/http"
)

func WriteResponse(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(message))
}
