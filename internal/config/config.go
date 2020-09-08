package config

import (
	"github.com/spf13/viper"
)

// CONFIG ...
var CONFIG Configuration

// Configuration ...
type Configuration struct {
	LoggerConfig LoggerConfig
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	// 日志文件名
	Filename string
	// 单个文件最大大小，单位: MB
	MaxSize int
	// 日志文件保存日期，单位: day
	MaxAge int
	// 日志文件保存数量
	MaxBackups int
	// 是否压缩日志文件
	Compress bool
	// 是否开启 debug 模式
	DebugMode bool
}

// InitConfiguration ...
func InitConfiguration(configName string, configPath string, config interface{}) error {
	vp := viper.New()
	vp.SetConfigName(configName)
	vp.AutomaticEnv()

	vp.AddConfigPath(configPath)

	if err := vp.ReadInConfig(); err != nil {
		return err
	}

	err := vp.Unmarshal(config)
	if err != nil {
		return err
	}

	return nil
}
