package cmd

import (
	"Panda/cmd/panda"

	"github.com/spf13/cobra"
)

var port string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server 模式",
	Long:  "Server 模式",
	Run: func(cmd *cobra.Command, args []string) {
		panda.Server(port)
	},
}

func init() {
	serverCmd.Flags().StringVar(&port, "port", "8080", "代理服务器端口号，默认 8080")
	rootCmd.AddCommand(serverCmd)
}
