package panda

import (
	"Panda/internal/core"
	"Panda/utils"
	"net"
)

// Client 是 Panda 的 Client 模式的实际入口
func Client(localePort string, remotePort string) {

	// proxy server 地址
	proxyServerAddr, err := net.ResolveTCPAddr("tcp", ":"+remotePort)
	if err != nil {
		utils.Logger.Error(err)
	}
	utils.Logger.Debug("连接远程服务器: ", remotePort+"....")

	// 监听本地
	listenAddr, err := net.ResolveTCPAddr("tcp", ":"+localePort)
	if err != nil {
		utils.Logger.Error(err)
	}
	utils.Logger.Debug("监听本地端口: ", localePort)

	listener, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		utils.Logger.Error(err)
	}

	for {
		client, err := listener.AcceptTCP()
		if err != nil {
			utils.Logger.Error(err)
		}
		go handleProxyRequest(client, proxyServerAddr)
	}
}

func handleProxyRequest(client *net.TCPConn, proxyServerAddr *net.TCPAddr) {

	// 连接 Proxy Server
	dstServer, err := net.DialTCP("tcp", nil, proxyServerAddr)
	if err != nil {
		utils.Logger.Debug("连接 Proxy Server 错误!!!")
		return
	}
	defer dstServer.Close()
	defer client.Close()

	core.SocksClient(client, dstServer)
}
