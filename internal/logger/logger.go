package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var ErrLog *log.Logger
var DebugLog *log.Logger
var WarnLog *log.Logger
var InfoLog *log.Logger

func init() {
	logFolder := os.Getenv("LOGS_FOLDER")
	err := os.MkdirAll(logFolder, 0666)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	logPath := filepath.Join(logFolder, fmt.Sprintf("%s.log", time.Now().Format("Jan-2-15:04:05")))
	f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}
	ErrLog = log.New(io.MultiWriter(os.Stderr, f), "[ERROR] ", log.LstdFlags|log.Lshortfile)
	if os.Getenv("DEBUG") == "TRUE" {
		DebugLog = log.New(io.MultiWriter(os.Stderr, f), "[DEBUG] ", log.LstdFlags|log.Lshortfile)
	} else {
		DebugLog = log.New(f, "[DEBUG] ", log.LstdFlags|log.Lshortfile)
	}
	WarnLog = log.New(io.MultiWriter(os.Stderr, f), "[WARN] ", log.LstdFlags|log.Lshortfile)
	InfoLog = log.New(io.MultiWriter(os.Stderr, f), "[INFO] ", log.LstdFlags)
}
