package internal

import (
	"net/http"

	"github.com/bennycio/bundle/internal/logger"
)

func HttpError(w http.ResponseWriter, err error, status int) {
	logger.ErrLog.Print(err.Error())
	http.Error(w, err.Error(), status)
}
