package panda

import (
	"Panda/internal/core"
	"Panda/utils"
	"net"
)

// Client 是 Panda 的 Client 模式的实际入口
func Client(localePort string, remotePort string, method int) {

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
		go handleProxyRequest(client, proxyServerAddr, toMethod(method))
	}
}

func handleProxyRequest(client *net.TCPConn, proxyServerAddr *net.TCPAddr, method byte) {

	// 连接 Proxy Server
	dstServer, err := net.DialTCP("tcp", nil, proxyServerAddr)
	if err != nil {
		utils.Logger.Debug("连接 Proxy Server 错误!!!")
		return
	}

	core.SocksClient(client, dstServer, method)
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
