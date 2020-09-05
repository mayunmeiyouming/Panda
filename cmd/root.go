package cmd

import (
	"Panda/cmd/panda"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listen string
var socks string

var rootCmd = &cobra.Command{
	Use:   "panda",
	Short: "Panda 是一个由 Golang 实现的代理服务器",
	Long: "Panda 是一个由 Golang 实现的代理服务器",
	Run: func(cmd *cobra.Command, args []string) {
		panda.Server(listen, socks)
	},
}

func init() {
	rootCmd.Flags().StringVar(&listen, "listen", "127.0.0.1:8080", "服务器端地址 address:port")
	rootCmd.Flags().StringVar(&socks, "socks", "127.0.0.1:1080", "本地代理地址，socks address:port")
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
