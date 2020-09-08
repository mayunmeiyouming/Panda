package core

import "net"

// SocksServe ...
func SocksServe(conn *net.TCPConn) {
	// 协商阶段
	SocksAuth(conn)
}
