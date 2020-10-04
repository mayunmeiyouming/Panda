package panda

import (
	"Panda/pkg/core"
	"Panda/utils"
	"net"
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
		go tcpRemote(addr, ciph.StreamConn)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}

func tcpRemote(addr string, shadow func(net.Conn) net.Conn) {
	listenAddr, err := net.ResolveTCPAddr("tcp", addr)
	l, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		utils.Logger.Fatal("监听端口失败，端口可能被占用")
	}

	for {
		utils.Logger.Info("等待连接")
		client, err := l.AcceptTCP()
		if err != nil {
			utils.Logger.Error(err)
		}
		utils.Logger.Debug("正在处理请求中")
		go func() {
			defer client.Close()

			sc := shadow(client)
			core.SocksServe(sc)
		}()
	}
}
