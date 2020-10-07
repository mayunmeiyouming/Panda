package panda

import (
	"Panda/pkg/core"
	"Panda/pkg/proxy"
	"Panda/utils"
	"os"
	"os/signal"
	"syscall"
)

// Client 是 Panda 的 Client 模式的实际入口
func Client(http string, socks string, remoteAddr string, cipher string, password string) {
	ciph, err := core.PickCipher(cipher, password)
	if err != nil {
		utils.Logger.Fatal("获取加密 cipher 失败: ", err)
	}

	if socks != "" {
		go proxy.SocksLocal(socks, remoteAddr, ciph.StreamConn)
	}

	if http != "" {
		go proxy.HTTPLocal(http, remoteAddr, ciph.StreamConn)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
