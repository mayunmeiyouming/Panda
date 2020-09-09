package core

import (
	"io"
	"net"
	"sync"
)

// SocksServe ...
func SocksServe(conn *net.TCPConn) {
	// 协商阶段
	socksAddressRequest, err := SocksAuth(conn)
	if err != nil {
		return
	}

	// 代理阶段
	destinationServer, err := net.DialTCP("tcp", nil, socksAddressRequest.RAWADDR)
	if err != nil {
		return
	}
	defer destinationServer.Close()

	wg := new(sync.WaitGroup)
	wg.Add(2)

	// 本地的内容copy到远程端
	go func() {
		defer wg.Done()
		io.Copy(destinationServer, conn)
	}()

	// 远程得到的内容copy到源地址
	go func() {
		defer wg.Done()
		io.Copy(conn, destinationServer)
	}()
	wg.Wait()
	conn.Close()

}
