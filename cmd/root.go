package cmd

import (
	"Panda/cmd/panda"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "panda",
	Short: "Panda 是一个由 Golang 实现的代理服务器",
	Long: "Panda 是一个由 Golang 实现的代理服务器",
	Run: func(cmd *cobra.Command, args []string) {
		panda.Server()
	},
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
