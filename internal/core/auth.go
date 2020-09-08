package core

import (
	"Panda/utils"
	"net"
)

// SocksAuthRequest ...
type SocksAuthRequest struct {
	VERSION  int32
	NMETHODS int32
	METHODS  int32
}

// SocksAuth ...
func SocksAuth(conn *net.TCPConn) {
	socksAuthRequest, n := parseSocksAuthRequest(conn)
	responseAuth(conn, socksAuthRequest, n)
	buf := make([]byte, 1024)
	n, _ = conn.Read(buf[0:])
	if n != 0 {
		utils.Log.Debug("success")
	} else {
		utils.Log.Warn("fail")
	}

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
	// conn.Write(b)
	utils.Log.Debug(conn.RemoteAddr())
	// address := conn.RemoteAddr()
	//获得了请求的host和port，就开始拨号吧
	// server, err := net.Dial("tcp", address.String())
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// server.Write(b[:])
	if len == 3 {
		// io.Copy(conn, bytes.NewReader(b))
		utils.Log.Debug(b[0:2])
		c, err := conn.Write(b[0:2])
		utils.Log.Debug(c, err)
	} else {
		return
	}
}
