package logger

import (
	"log"
	"os"
)

var ErrLog = log.New(os.Stderr, "Err: ", log.LstdFlags|log.Lshortfile)
