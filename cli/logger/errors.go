package logger

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

var ErrLog *log.Logger

func init() {
	if viper.GetBool("debug") {
		ErrLog = log.New(os.Stderr, "[ERROR]", log.LstdFlags|log.Lshortfile)
	} else {
		ErrLog = log.New(os.Stderr, "Err: ", log.Lmsgprefix)
	}

}
