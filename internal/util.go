package internal

import (
	"encoding/json"
	"net/http"
)

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func WriteResponse(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(message))
}
