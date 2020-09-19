package utils

import (
	"Panda/internal/config"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/sirupsen/logrus"
)

// Logger ...
var Logger = logrus.New()

// InitLogger 初始化
func InitLogger() {
	// 设置 log 的级别为 Info
	Logger.SetLevel(logrus.InfoLevel)
	if config.CONFIG.LoggerConfig.DebugMode {
		Logger.SetLevel(logrus.TraceLevel)
		Logger.SetReportCaller(true)
	}

	Logger.SetOutput(os.Stdout)

	// 设置日志的 json 格式
	Logger.Formatter = &logrus.JSONFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Dir(f.File) + "/" + path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}
}
