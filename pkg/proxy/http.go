package proxy

import (
	"Panda/pkg/socks"
	"Panda/utils"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
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

// 构造地址请求
func makeAddrRequest(host string, port string) []byte {
	addr := make([]byte, 0)

	address := net.ParseIP(host)
	if address != nil {
		// IPv4
		if len(address) == 4 {
			addr = append(addr, 0x01)
		} else {
			// IPv6
			addr = append(addr, 0x04)
		}
	} else {
		addr = append(addr, 0x03)
		addr = append(addr, byte(len([]byte(host)))) // 域名字节长度
	}

	// 域名
	addr = append(addr, []byte(host)...)

	// 端口
	b := []byte{0, 0}
	r, _ := strconv.Atoi(port)
	binary.BigEndian.PutUint16(b, uint16(r))
	addr = append(addr, b[:2]...)

	return addr
}
