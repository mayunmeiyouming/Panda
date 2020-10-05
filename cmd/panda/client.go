package panda

import (
	"os"
	"os/signal"
	"syscall"
	"Panda/pkg/proxy"
)

// Client 是 Panda 的 Client 模式的实际入口
func Client(http string, socks string, remoteAddr string, cipher string, password string) {
	if socks != "" {

	}

	if http != "" {
		go proxy.HTTPLocal(http, remoteAddr, cipher, password)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}
