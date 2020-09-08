package panda

import (
	"Panda/core"
	"log"
	"net"
)

// Server 是 Panda 的实际入口
func Server(port string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	listenAddr, err := net.ResolveTCPAddr("tcp", ":"+port)
	l, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		log.Panic(err)
	}
	log.Println("等待连接")
	for {
		client, err := l.AcceptTCP()
		if err != nil {
			log.Panic(err)
		}
		log.Println("正在处理请求中")
		go handleClientRequest(client)
	}
}

func handleClientRequest(client *net.TCPConn) {
	if client == nil {
		return
	}
	defer client.Close()
	core.SocksAuth(client)
}
