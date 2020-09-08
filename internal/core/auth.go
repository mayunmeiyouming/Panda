package core

import (
	"Panda/utils"
	"encoding/binary"
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
	VER       uint8
	CMD       uint8
	RSV       uint8 // 保留位
	ATYP      uint8
	DSTADDR   []byte
	DSTPORT   uint16
	DSTDOMAIN string
}

// SocksAuth 是协商认证阶段
func SocksAuth(conn *net.TCPConn) error {
	// 协商认证方法
	socksAuthRequest, err := parseSocksAuthRequest(conn)
	if err != nil {
		return err
	}
	err = responseAuth(conn, socksAuthRequest)
	if err != nil {
		return err
	}

	// 获取代理地址
	parseSocksAddressRequest(conn)

	return nil
}

func parseSocksAuthRequest(conn *net.TCPConn) (*SocksAuthRequest, error) {
	b := make([]byte, 1024)
	n, err := conn.Read(b)

	if n >= 3 {
		socksAuthRequest := &SocksAuthRequest{
			VERSION:  int32(b[0]),
			NMETHODS: int32(b[1]),
			METHODS:  int32(b[2]),
		}
		utils.Log.Debug(socksAuthRequest)
		return socksAuthRequest, nil
	}
	utils.Log.Debug("认证协议格式错误")
	utils.Log.Debug("length: ", n)
	return nil, err
}

func responseAuth(conn *net.TCPConn, socks *SocksAuthRequest) error {
	b := []byte{0x05, 0x00}
	utils.Log.Debug(conn.RemoteAddr())
	_, err := conn.Write(b[0:2])
	if err != nil {
		return err
	}
	return nil
}

func parseSocksAddressRequest(conn *net.TCPConn) (*SocksAddressRequest, error) {
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf[0:])

	socksAddressRequest := SocksAddressRequest{
		VER:  uint8(buf[0]),
		CMD:  uint8(buf[1]),
		RSV:  uint8(buf[2]),
		ATYP: uint8(buf[3]),
	}

	if socksAddressRequest.ATYP == 1 {
		// IPv4
		socksAddressRequest.DSTADDR = buf[4:8]
		utils.Log.Debug("IPv4")
	} else if socksAddressRequest.ATYP == 3 {
		// Domain
		socksAddressRequest.DSTDOMAIN = string(buf[5 : n-2])
		ipAddr, err := net.ResolveIPAddr("ip", socksAddressRequest.DSTDOMAIN)
		if err != nil {
			return nil, err
		}
		socksAddressRequest.DSTADDR = ipAddr.IP
		utils.Log.Debug("Domain")
	} else if socksAddressRequest.ATYP == 4 {
		// IPv6
		socksAddressRequest.DSTADDR = buf[4 : 4+net.IPv6len]
		utils.Log.Debug("IPv6")
	}

	socksAddressRequest.DSTPORT = binary.BigEndian.Uint16(buf[n-2 : n])

	utils.Log.Debug(socksAddressRequest)
	return &socksAddressRequest, nil
}
