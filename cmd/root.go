package cmd

import (
	"Panda/internal/config"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"Panda/utils"
)

var rootCmd = &cobra.Command{
	Use:   "panda",
	Short: "Panda 是一个由 Golang 实现的代理服务器",
	Long:  "Panda 是一个由 Golang 实现的代理服务器",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initResource()
	},
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initResource()  {
	config.InitConfiguration("config", "./configs/", &config.CONFIG)
	utils.InitLogger()

	if config.CONFIG.LoggerConfig.DebugMode {
		utils.Logger.Info("running on debug mode")
	}
}

