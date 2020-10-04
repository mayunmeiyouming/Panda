package panda

import (
	"Panda/pkg/core"
	"Panda/utils"
	"net"
	"os"
	"os/signal"
	"syscall"
)

// Client 是 Panda 的 Client 模式的实际入口
func Client(localAddr string, remoteAddr string, cipher string, password string) {
	go socksLocal(localAddr, remoteAddr, cipher, password)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}

func socksLocal(localAddr string, remoteAddr string, cipher string, password string) {
	// proxy server 地址
	proxyServerAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		utils.Logger.Error(err)
	}
	utils.Logger.Debug("连接远程服务器: ", remoteAddr+"....")

	// 监听本地
	listenAddr, err := net.ResolveTCPAddr("tcp", localAddr)
	if err != nil {
		utils.Logger.Error(err)
	}
	utils.Logger.Debug("监听本地端口: ", localAddr)

	listener, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		utils.Logger.Error(err)
	}

	ciph, err := core.PickCipher(cipher, password)
	if err != nil {
		utils.Logger.Fatal(err)
	}

	for {
		client, err := listener.AcceptTCP()
		if err != nil {
			utils.Logger.Error(err)
		}
		defer client.Close()

		go handleProxyRequest(client, proxyServerAddr, ciph.StreamConn)
	}
}

func handleProxyRequest(client *net.TCPConn, proxyServerAddr *net.TCPAddr, shadow func(net.Conn) net.Conn) {

	// 连接 Proxy Server
	dstServer, err := net.DialTCP("tcp", nil, proxyServerAddr)
	if err != nil {
		utils.Logger.Debug("连接 Proxy Server 错误!!!, ", err)
		return
	}
	defer dstServer.Close()

	sd := shadow(dstServer)
	core.SocksClient(client, sd)
}

func toMethod(method int) byte {
	var res byte
	switch method {
	case 0:
		res = 0x00
	case 1:
		res = 0x80
	default:
		utils.Logger.Fatal("请选择正确的加密方式")
	}

	return res
}
