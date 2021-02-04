package common

import (
	"io"
	"os"

	logging "github.com/sirupsen/logrus"
)

// Logger a logrus logger
var Logger = logging.New()

// SetupLogger setup log and return logger
func SetupLogger(debug bool, logFileName string) (*logging.Logger, error) {
	logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return Logger, InitailizeError(err.Error())
	}

	Logger.SetFormatter(&logging.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	Logger.SetOutput(io.MultiWriter(os.Stdout, logFile))
	if debug {
		Logger.SetLevel(logging.DebugLevel)
	} else {
		Logger.SetLevel(logging.InfoLevel)
	}

	return Logger, nil
}
