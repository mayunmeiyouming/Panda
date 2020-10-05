package panda

import (
	"Panda/pkg/core"
	"Panda/pkg/proxy"
	"Panda/utils"
	"os"
	"os/signal"
	"syscall"
)

// Server 是 Panda 的 Server 模式的实际入口
func Server(addr string, tcp bool, cipher string, password string) {
	ciph, err := core.PickCipher(cipher, password)
	if err != nil {
		utils.Logger.Fatal(err)
	}

	if tcp {
		utils.Logger.Info("开启 TCP Listen")
		go proxy.TCPRemote(addr, ciph.StreamConn)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
