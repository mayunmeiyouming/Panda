package cmd

import (
	"Panda/cmd/panda"

	"github.com/spf13/cobra"
)

var localePort string
var remotePort string
var method int // 加密算法

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Client 模式",
	Long:  "Client 模式",
	Run: func(cmd *cobra.Command, args []string) {
		panda.Client(localePort, remotePort, method)
	},
}

func init() {
	clientCmd.Flags().StringVar(&localePort, "localePort", "2080", "本地监听端口号，默认 2080")
	clientCmd.Flags().StringVar(&remotePort, "remotePort", "8080", "代理服务器端口号，默认 8080")
	clientCmd.Flags().IntVar(&method, "method", 0, "代理服务器加密算法，默认无加密")
	rootCmd.AddCommand(clientCmd)
}
