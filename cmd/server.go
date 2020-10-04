package cmd

import (
	"Panda/cmd/panda"
	"Panda/pkg/core"
	"Panda/utils"
	"strings"

	"github.com/spf13/cobra"
)



var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Server 模式",
	Long:  "Server 模式",
	PreRun: func(cmd *cobra.Command, args []string) {
		var err error
		if strings.HasPrefix(addr, "ss://") {
			addr, cipher, password, err = parseURL(addr)
			if err != nil {
				utils.Logger.Fatal(err)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logger.Info(cipher)
		panda.Server(addr, tcp, cipher, password)
	},
}

func init() {
	serverCmd.Flags().StringVar(&addr, "addr", "ss://AES-256-GCM:astaxie12798akljzmknmfahkjkljlfk@:8080", "server listen address or url")
	serverCmd.Flags().StringVar(&cipher, "cipher", "DUMMY", "available ciphers: "+strings.Join(core.ListCipher(), " "))
	serverCmd.Flags().StringVar(&password, "password", "astaxie12798akljzmknmfahkjkljlfk", "password")
	serverCmd.Flags().StringVar(&plugin, "plugin", "", "Enable SIP003 plugin. (e.g., v2ray-plugin)")
	serverCmd.Flags().StringVar(&pluginOpts, "plugin-opts", "", "Set SIP003 plugin options. (e.g., \"server;tls;host=mydomain.me\")")
	serverCmd.Flags().BoolVar(&tcp, "tcp", true, "enable TCP support")
	serverCmd.Flags().BoolVar(&udp, "udp", false, "enable UDP support")
	rootCmd.AddCommand(serverCmd)
}
