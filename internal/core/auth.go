package core

import (
	"Panda/utils"
	"net"
)

// SocksAuthRequest 是客户端的协商请求
type SocksAuthRequest struct {
	VERSION  int32
	NMETHODS int32
	METHODS  int32
}

// SocksAddressRequest 是客户端告诉服务端目标地址的请求
type SocksAddressRequest struct {
	VERSION                int32
	COMMAND                int32
	RSV                    int32 // 保留位
	AddressType            int32
	DestinationAddress     int32
	DestinationAddressPORT int32
}

// SocksAuth ...
func SocksAuth(conn *net.TCPConn) {
	// 协商认证方法
	socksAuthRequest, n := parseSocksAuthRequest(conn)
	responseAuth(conn, socksAuthRequest, n)

	// 获取代理地址
	getSocksAddress(conn)

}

func parseSocksAuthRequest(conn *net.TCPConn) (*SocksAuthRequest, int) {
	b := make([]byte, 1024)
	n, _ := conn.Read(b)

	if n >= 3 {
		socksAuthRequest := &SocksAuthRequest{
			VERSION:  int32(b[0]),
			NMETHODS: int32(b[1]),
			METHODS:  int32(b[2]),
		}
		utils.Log.Debug(socksAuthRequest)
		return socksAuthRequest, n
	}
	utils.Log.Debug("认证协议格式错误")
	utils.Log.Debug("length: ", n)
	return nil, n
}

func responseAuth(conn *net.TCPConn, socks *SocksAuthRequest, len int) {
	b := []byte{0x05, 0x00}
	utils.Log.Debug(conn.RemoteAddr())
	if len >= 3 {
		utils.Log.Debug(b[0:2])
		c, err := conn.Write(b[0:2])
		utils.Log.Debug(c, err)
	} else {
		return
	}
}

func getSocksAddress(conn *net.TCPConn) {
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf[0:])
	utils.Log.Debug("length: ", n)
}
