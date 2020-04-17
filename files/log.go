package files

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

// SetupLogLevel gets the debug level from the configuration and sets it
func SetupLogLevel(configuration Configuration) {
	level, err := log.ParseLevel(configuration.DebugLevel)
	if err != nil {
		log.Fatal(err)
	}

	log.SetLevel(level)
}

// SetupLog creates the log file with a given filename
func SetupLog(filename string) {
	// Create the log file if doesn't exist. And append to it if it already exists.
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	Formatter := new(log.TextFormatter)

	// You can change the Timestamp format. But you have to use the same date and time.
	// "2006-02-02 15:04:06" Works. If you change any digit, it won't work
	// ie "Mon Jan 2 15:04:05 MST 2006" is the reference time. You can't change it
	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)

	if err != nil {
		// Cannot open log file. Logging to stderr
		fmt.Println(err)
	} else {
		log.SetOutput(f)
	}
}
