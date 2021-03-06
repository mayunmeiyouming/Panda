package socks

import (
	"Panda/utils"
	"encoding/binary"
	"errors"
	"net"
	"strconv"
	"time"
)

// SocksAuthRequest 是客户端的协商请求
type SocksAuthRequest struct {
	VERSION  int32
	NMETHODS int32
	METHODS  []byte
}

// SocksAddressRequest 是客户端告诉服务端目标地址的请求
type SocksAddressRequest struct {
	VER       uint8
	CMD       uint8
	RSV       uint8 // 保留位
	ATYP      uint8
	ADDR      string
	PORT      string
	DSTADDR   []byte
	DSTPORT   uint16
	DSTDOMAIN string
	RAWADDR   *net.TCPAddr
}

// SocksAuth 是协商认证阶段
func SocksAuth(conn net.Conn) (*SocksAddressRequest, *byte, error) {
	// 协商认证方法
	socksAuthRequest, err := parseSocksAuthRequest(conn)
	if err != nil {
		return nil, nil, err
	}
	err = responseAuth(conn, socksAuthRequest)
	if err != nil {
		return nil, nil, err
	}

	// 获取代理地址
	socksAddressRequest, err := parseSocksAddressRequest(conn)
	if err != nil {
		return nil, nil, err
	}
	// 服务端回复
	responseSocksAddressRequest(conn, socksAddressRequest)

	return socksAddressRequest, &socksAuthRequest.METHODS[0], nil
}

// 第一阶段协商，解析 client 的包
func parseSocksAuthRequest(conn net.Conn) (*SocksAuthRequest, error) {
	b := make([]byte, 257)
	n, err := conn.Read(b)

	if n >= 3 {
		socksAuthRequest := &SocksAuthRequest{
			VERSION:  int32(b[0]),
			NMETHODS: int32(b[1]),
			METHODS:  b[2:n],
		}
		utils.Logger.Debug("一次协商请求: ", *socksAuthRequest)
		if socksAuthRequest.VERSION != 5 {
			return nil, errors.New("只支持 SOCKS5 协议，该协议为 SOCKS"+strconv.Itoa(int(socksAuthRequest.VERSION)))
		}
		return socksAuthRequest, nil
	}
	utils.Logger.Debug("认证协议格式错误")
	utils.Logger.Debug("length: ", n)
	return nil, err
}

func responseAuth(conn net.Conn, socks *SocksAuthRequest) error {
	b := []byte{0x05, 0x00}
	utils.Logger.Debug("一次协商回复: ", b[0:2])
	_, err := conn.Write(b[0:2])
	if err != nil {
		return err
	}
	return nil
}

// 解析第二阶段协商的 client 的包
func parseSocksAddressRequest(conn net.Conn) (*SocksAddressRequest, error) {
	buf := make([]byte, 2048)
	n, _ := conn.Read(buf[0:])

	// 解决网络延时问题
	if n == 0 {
		for i := 0; i < 5; i++ {
			time.Sleep(50 * time.Microsecond)
			n, _ = conn.Read(buf[0:])
			if n != 0 {
				break
			}
		}
		if n == 0 {
			utils.Logger.Error("error")
			return nil, errors.New("可能网络不通")
		}

	}

	// 判断第二阶段请求包是否完整
	if len(buf) < 10 {
		return nil, errors.New("第二阶段请求包不完整")
	}

	socksAddressRequest := SocksAddressRequest{
		VER:  uint8(buf[0]),
		CMD:  uint8(buf[1]),
		RSV:  uint8(buf[2]),
		ATYP: uint8(buf[3]),
	}

	if socksAddressRequest.ATYP == 1 {
		// IPv4
		socksAddressRequest.DSTADDR = buf[4:8]
		socksAddressRequest.ADDR = string(buf[4:8])
		// utils.Logger.Debug("IPv4")
	} else if socksAddressRequest.ATYP == 3 {
		// Domain
		l := int(buf[4])
		if l+7 != n {
			return nil, errors.New("第二阶段请求包不完整")
		}

		socksAddressRequest.DSTDOMAIN = string(buf[5 : n-2])
		socksAddressRequest.ADDR = string(buf[5 : n-2])
		ipAddr, err := net.ResolveIPAddr("ip", socksAddressRequest.DSTDOMAIN)
		if err != nil {
			return nil, err
		}
		socksAddressRequest.DSTADDR = ipAddr.IP[len(ipAddr.IP)-4:]
		utils.Logger.Debug("Domain")
	} else if socksAddressRequest.ATYP == 4 {
		// IPv6
		socksAddressRequest.DSTADDR = buf[4 : 4+net.IPv6len]
		socksAddressRequest.ADDR = string(buf[4 : 4+net.IPv6len])
		utils.Logger.Debug("IPv6")
	}

	socksAddressRequest.DSTPORT = binary.BigEndian.Uint16(buf[n-2 : n])

	socksAddressRequest.RAWADDR = &net.TCPAddr{
		IP:   socksAddressRequest.DSTADDR,
		Port: int(socksAddressRequest.DSTPORT),
	}

	socksAddressRequest.PORT = strconv.Itoa(int(socksAddressRequest.DSTPORT))

	utils.Logger.Debug("二次协商请求: ", socksAddressRequest)
	utils.Logger.Debug("客户端需要访问的服务器地址: ", socksAddressRequest.DSTADDR)
	return &socksAddressRequest, nil
}

func responseSocksAddressRequest(conn net.Conn, socks *SocksAddressRequest) error {
	response := []byte{0x05, 0x00, 0x00}

	response = append(response, socks.ATYP)
	if socks.ATYP == 1 {
		// IPv4
		response = append(response, socks.DSTADDR...)
	} else if socks.ATYP == 3 {
		// Domain
		b := []byte(socks.DSTDOMAIN)
		response = append(response, byte(len(b)))
		response = append(response, b...)
	} else if socks.ATYP == 4 {
		// IPv6
		response = append(response, socks.DSTADDR...)
	}

	b := []byte{0, 0} // 初始化
	binary.BigEndian.PutUint16(b, socks.DSTPORT)
	response = append(response, b[:2]...)

	utils.Logger.Debug("二次协商回复: ", response)
	conn.Write(response)

	return nil
}
