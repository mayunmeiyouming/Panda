package cmd

import (
	"Panda/cmd/panda"
	"Panda/pkg/core"
	"Panda/utils"
	"strings"

	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Client 模式",
	Long:  "Client 模式",
	PreRun: func(cmd *cobra.Command, args []string) {
		var err error
		if strings.HasPrefix(socks, "ss://") {
			remoteAddr, cipher, password, err = parseURL(socks)
			if err != nil {
				utils.Logger.Fatal(err)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logger.Info("加密: ", cipher)
		panda.Client(http, socks, remoteAddr, cipher, password)
	},
}

func init() {
	clientCmd.Flags().StringVar(&socks, "s", ":2080", "server listen address or url")
	clientCmd.Flags().StringVar(&http, "http", "", "本地监听端口号，默认 2080")
	clientCmd.Flags().StringVar(&remoteAddr, "remoteAddr", ":8080", "代理服务器端口号，默认 8080")
	clientCmd.Flags().StringVar(&cipher, "cipher", "AES_256_GCM", "available ciphers: "+strings.Join(core.ListCipher(), " "))
	clientCmd.Flags().StringVar(&password, "password", "astaxie12798akljzmknmfahkjkljlfk", "password")
	clientCmd.Flags().BoolVar(&udpsocks, "u", false, "(client-only) Enable UDP support for SOCKS")
	rootCmd.AddCommand(clientCmd)
}
