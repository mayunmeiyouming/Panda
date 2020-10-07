package proxy

import (
	"Panda/pkg/socks"
	"Panda/utils"
	"fmt"
	"net"
)

// HTTPLocal is only client
func HTTPLocal(httpAddr string, remoteAddr string, shadow func(net.Conn) net.Conn) {
	// proxy server 地址
	proxyServerAddr, err := net.ResolveTCPAddr("tcp", remoteAddr)
	if err != nil {
		utils.Logger.Error(err)
	}
	utils.Logger.Debug("连接远程服务器: ", remoteAddr+"....")

	// 监听本地
	listenAddr, err := net.ResolveTCPAddr("tcp", httpAddr)
	if err != nil {
		utils.Logger.Error(err)
	}
	utils.Logger.Debug("监听本地端口: ", httpAddr)

	listener, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		utils.Logger.Error(err)
	}

	for {
		client, err := listener.AcceptTCP()
		if err != nil {
			utils.Logger.Error(err)
		}
		defer client.Close()

		go handleProxyRequest(client, proxyServerAddr, shadow)
	}
}

func handleProxyRequest(client *net.TCPConn, proxyServerAddr *net.TCPAddr, shadow func(net.Conn) net.Conn) {

	// 连接 Proxy Server
	dstServer, err := net.DialTCP("tcp", nil, proxyServerAddr)
	if err != nil {
		utils.Logger.Error("连接 Proxy Server 错误!!!, ", err)
		return
	}
	defer dstServer.Close()

	sd := shadow(dstServer)

	buff, host, port, err := socks.ParseHTTP(client)
	if err != nil {
		return
	}

	addr := makeAddrRequest(*host, *port)

	sd.Write(addr)

	// 转发消息
	if *port == "443" {
		fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		sd.Write(buff)
		// utils.Logger.Debug("HTTP发送成功")
	}

	//进行转发
	relay(sd, client)
}
