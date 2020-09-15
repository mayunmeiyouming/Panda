package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Panic(err)
	}
	for {
		client, err := l.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleClientRequest(client)
	}
}
func handleClientRequest(client net.Conn) {
	if client == nil {
		return
	}
	defer client.Close()
	b := make([]byte, 1024)
	n, err := client.Read(b[:])
	if err != nil {
		log.Println(err)
		return
	}
	var address string
	re := bytes.NewReader(b[:n])
	a := bufio.NewReader(re)
	request, err := http.ReadRequest(a)
	if err != nil {
		log.Println(err)
		return
	}

	method := request.Method

	dst := strings.Split(request.Host, ":")
	var dstPort = "80"
	if len(dst) == 2 && dst[1] != "" {
		dstPort = dst[1]
	}

	if dstPort == "80" {
		address = request.Host + ":" + dstPort
	} else {
		address = request.Host
	}

	log.Println(address)
	server, err := net.Dial("tcp", address)
	if err != nil {
		log.Println(err)
		return
	}
	//如果使用https协议，需先向客户端表示连接建立完毕
	if method == "CONNECT" {
		fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		//如果使用http协议，需将从客户端得到的http请求转发给服务端
		server.Write(b[:n])
	}

	//将客户端的请求转发至服务端，将服务端的响应转发给客户端。io.Copy为阻塞函数，文件描述符不关闭就不停止
	go io.Copy(server, client)
	io.Copy(client, server)
}
