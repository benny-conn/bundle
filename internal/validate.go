package internal

import (
	"net/http"
	"os"
)

var validStatusCodes = []int{http.StatusOK, http.StatusAccepted, http.StatusCreated, http.StatusFound, http.StatusSeeOther, http.StatusProcessing, http.StatusContinue}

func IsValidPath(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func IsRespError(r *http.Response) bool {
	for _, v := range validStatusCodes {
		if r.StatusCode == v {
			return false
		}
	}
	return true
}
