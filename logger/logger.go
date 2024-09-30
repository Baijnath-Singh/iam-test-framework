package logger

import (
	"log"
	"os"
)

var logger *log.Logger

// NewLogger initializes a new logger instance.
func NewLogger(logFilePath string) {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger = log.New(file, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// LogError logs an error message.
func LogError(message string) {
	if logger != nil {
		logger.Println("ERROR: " + message)
	}
}
