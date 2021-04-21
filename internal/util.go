package internal

import (
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func WriteResponse(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	w.Write([]byte(message))
}

func NewCaseInsensitiveRegex(value string) primitive.Regex {
	return primitive.Regex{Pattern: value, Options: "i"}
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
