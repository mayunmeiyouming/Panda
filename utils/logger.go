package utils

import (
	"Panda/internal/config"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

// Logger  global logger
var Logger = logrus.New()

// InitLogger  set logger
func InitLogger() {
	// only log the InfoLevel or above
	Logger.SetLevel(logrus.InfoLevel)
	if config.CONFIG.LoggerConfig.DebugMode {
		// log all the level
		Logger.SetLevel(logrus.TraceLevel)

		// set whether print caller
		Logger.SetReportCaller(true)
	}

	// output to stdout instead of the default stderr
	// can be any io.Writer
	Logger.SetOutput(os.Stdout)

	// set json formatter
	Logger.Formatter = &logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Dir(f.File) + "/" + path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}
}
