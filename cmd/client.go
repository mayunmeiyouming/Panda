package cmd

import (
	"Panda/cmd/panda"

	"github.com/spf13/cobra"
)

var localePort string
var remotePort string

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Client 模式",
	Long:  "Client 模式",
	Run: func(cmd *cobra.Command, args []string) {
		panda.Client(localePort, remotePort)
	},
}

func init() {
	clientCmd.Flags().StringVar(&localePort, "localePort", "1080", "本地监听端口号，默认 1080")
	clientCmd.Flags().StringVar(&remotePort, "remotePort", "8080", "代理服务器端口号，默认 8080")
	rootCmd.AddCommand(clientCmd)
}
