package panda

import (
	"fmt"
	"log"
	"net"
	"runtime"
	"strings"
	"time"
)

var host = "127.0.0.1"
var remotePort = "8848"
var localPort = "80"

// 与 browser 相关的 conn
type browser struct {
	conn net.Conn
	er   chan bool
	writ chan bool
	recv chan []byte
	send chan []byte
}

// 读取 browser 过来的数据
func (br browser) read() {
	for {
		var recv []byte = make([]byte, 10240)
		n, err := br.conn.Read(recv)
		if err != nil {
			br.writ <- true
			br.er <- true
			//fmt.Println("读取browser失败", err)
			break
		}
		br.recv <- recv[:n]
	}
}

// 把数据发送给 browser
func (br browser) write() {
	for {
		var send []byte = make([]byte, 10240)
		select {
		case send = <-br.send:
			br.conn.Write(send)
		case <-br.writ:
			//fmt.Println("写入browser进程关闭")
			break
		}
	}
}

// 与 server 相关的 conn
type server struct {
	conn net.Conn
	er   chan bool
	writ chan bool
	recv chan []byte
	send chan []byte
}

// read 将读取服务器发送过来的数据
func (srv *server) read() {
	//isheart与timeout共同判断是不是自己设定的SetReadDeadline
	var isheart bool = false
	//20秒发一次心跳包
	srv.conn.SetReadDeadline(time.Now().Add(time.Second * 20))
	for {
		var recv []byte = make([]byte, 10240)
		n, err := srv.conn.Read(recv)
		if err != nil {
			if strings.Contains(err.Error(), "timeout") && !isheart {
				//fmt.Println("发送心跳包")
				srv.conn.Write([]byte("hh"))
				//4秒时间收心跳包
				srv.conn.SetReadDeadline(time.Now().Add(time.Second * 4))
				isheart = true
				continue
			}
			//浏览器有可能连接上不发消息就断开，此时就发一个0，为了与服务器一直有一条tcp通路
			srv.recv <- []byte("0")
			srv.er <- true
			srv.writ <- true
			//fmt.Println("没收到心跳包或者server关闭，关闭此条tcp", err)
			break
		}
		//收到心跳包
		if recv[0] == 'h' && recv[1] == 'h' {
			//fmt.Println("收到心跳包")
			srv.conn.SetReadDeadline(time.Now().Add(time.Second * 20))
			isheart = false
			continue
		}
		srv.recv <- recv[:n]
	}
}

//把数据发送给server
func (srv server) write() {

	for {
		var send []byte = make([]byte, 10240)

		select {
		case send = <-srv.send:
			srv.conn.Write(send)
		case <-srv.writ:
			//fmt.Println("写入server进程关闭")
			break
		}

	}

}

// Server 是 Panda 的实际入口
func Server() {
	target := net.JoinHostPort(host, remotePort)
	for {
		// 链接端口
		serverconn := dail(target)
		recv := make(chan []byte)
		send := make(chan []byte)
		// 1个位置是为了防止两个读取线程一个退出后另一个永远卡住
		er := make(chan bool, 1)
		writ := make(chan bool)
		next := make(chan bool)
		server := &server{serverconn, er, writ, recv, send}
		go server.read()
		go server.write()
		go handle(server, next)
		<-next
	}
}

//链接端口
func dail(target string) net.Conn {
	conn, err := net.Dial("tcp", target)
	logExit(err)
	return conn
}

//显示错误并退出
func logExit(err error) {
	if err != nil {
		log.Printf("出现错误，退出线程： %v\n", err)
		runtime.Goexit()
	}
}

//两个socket衔接相关处理
func handle(server *server, next chan bool) {
	var serverrecv = make([]byte, 10240)
	//阻塞这里等待server传来数据再链接browser
	fmt.Println("等待server发来消息")
	serverrecv = <-server.recv
	//连接上，下一个tcp连上服务器
	next <- true
	//fmt.Println("开始新的tcp链接，发来的消息是：", string(serverrecv))
	var browse *browser
	//server发来数据，链接本地80端口
	serverconn := dail("127.0.0.1:" + localPort)
	recv := make(chan []byte)
	send := make(chan []byte)
	er := make(chan bool, 1)
	writ := make(chan bool)
	browse = &browser{serverconn, er, writ, recv, send}
	go browse.read()
	go browse.write()
	browse.send <- serverrecv

	for {
		var serverrecv = make([]byte, 10240)
		var browserrecv = make([]byte, 10240)
		select {
		case serverrecv = <-server.recv:
			if serverrecv[0] != '0' {

				browse.send <- serverrecv
			}

		case browserrecv = <-browse.recv:
			server.send <- browserrecv
		case <-server.er:
			//fmt.Println("server关闭了，关闭server与browse")
			server.conn.Close()
			browse.conn.Close()
			runtime.Goexit()
		case <-browse.er:
			//fmt.Println("browse关闭了，关闭server与browse")
			server.conn.Close()
			browse.conn.Close()
			runtime.Goexit()
		}
	}
}
