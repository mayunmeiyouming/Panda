package core

import (
	"log"
	"net"
)

// SocksClientConnRequest ...
type SocksClientConnRequest struct {
	VERSION  int32
	NMETHODS int32
	METHODS  int32
}

// BuildConn ...
func BuildConn(conn net.Conn) {
	for {
		n := parseSocksClientRequest(conn)
		b := []byte{byte(5), 0x01}
		log.Println(b)
		// conn.Write(b)
		log.Println(conn.RemoteAddr())
		// address := conn.RemoteAddr()
		//获得了请求的host和port，就开始拨号吧
		// server, err := net.Dial("tcp", address.String())
		// if err != nil {
		// 	log.Println(err)
		// 	return
		// }
		// server.Write(b[:])
		if n == 3 {
			// io.Copy(conn, bytes.NewReader(b))
			log.Println(b[0:2])
			conn.Write(b[0:2])
		} else {
			return
		}
	}

}

func parseSocksClientRequest(conn net.Conn) int {
	b := make([]byte, 1024)
	n, _ := conn.Read(b)

	if n == 3 {
		socksClientConnRequest := &SocksClientConnRequest{
			VERSION:  int32(b[0]),
			NMETHODS: int32(b[1]),
			METHODS:  int32(b[2]),
		}
		log.Println(socksClientConnRequest)
		return n
	}
	log.Println("length: ", n)
	return n
}
