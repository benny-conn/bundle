package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

var errFile *os.File

var ErrLog = log.New(io.MultiWriter(os.Stderr, errFile), "Err: ", log.LstdFlags|log.Lshortfile)

func init() {
	f, err := os.OpenFile(os.Getenv("ERR_LOG"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err opening log file: %s\n", err.Error())
		return
	}
	errFile = f
}
