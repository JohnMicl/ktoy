package logger

import (
	"io"
	"ktoy/utils"
	"os"

	log "github.com/sirupsen/logrus"
)

func Loginit() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	// log.SetOutput(os.Stdout)
	logname := utils.GetLogFileName()
	file, err := os.OpenFile("release/x86_64/linux/logs/"+logname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	writers := []io.Writer{
		file,
		os.Stdout}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	if err == nil {
		log.SetOutput(fileAndStdoutWriter)
	} else {
		log.Panic("failed to log to file.")
	}

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	// Set reporter
	log.SetReportCaller(true)
}
