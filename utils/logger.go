package utils

import (
	"Panda/internal/config"

	"log"
)

// Log ...
var Log Logger

// Logger ...
type Logger struct {
	config config.LoggerConfig
}

// InitLogger ...
func InitLogger(loggerConfig config.LoggerConfig) {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	Log.config = loggerConfig
}

// Debug ...
func (logger Logger) Debug(v ...interface{}) {
	if logger.config.DebugMode == true {
		log.Println(v...)
	}
}

// Error ...
func (logger Logger) Error(v ...interface{}) {
	log.Println(v...)
}

// Warn ...
func (logger Logger) Warn(v ...interface{}) {
	log.Println(v...)
}
