package panda

import (
	"Panda/core"
	"log"
	"net"
)

// Server 是 Panda 的实际入口
func Server(port string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Panic(err)
	}
	log.Println("等待连接")
	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		log.Println("正在处理请求中")
		go handleClientRequest(client)
	}
}

func handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()
	core.BuildConn(client)
}
