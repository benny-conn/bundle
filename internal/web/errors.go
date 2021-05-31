package web

import (
	"errors"
	"net/http"
	"strings"

	"github.com/bennycio/bundle/internal/logger"
)

type errorData struct {
	Code    int
	Message string
}

func cleanError(err error) error {
	return errors.New(strings.ReplaceAll(strings.ReplaceAll(err.Error(), "rpc error: code = Unknown desc = ", ""), "mongo:", ""))
}
func handleError(w http.ResponseWriter, err error, code int) {
	msg := cleanError(err).Error()
	errData := errorData{
		Code:    code,
		Message: msg,
	}
	data := templateData{
		Error: errData,
	}
	e := tpl.ExecuteTemplate(w, "error", data)
	if e != nil {
		logger.ErrLog.Print(e)
	}
}
