package core

import (
	"Panda/utils"
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

	utils.Logger.Info("代理成功")
}

// SocksClient ...
func SocksClient(client *net.TCPConn, dstServer *net.TCPConn) {

	// socket5请求认证协商
	// 第一阶段协议版本及认证方式
	socksClientAuthResponse, err := RequestVersionAndMethodAuth(dstServer)
	if err != nil {
		return
	}

	// 第二阶段根据认证方式执行对应的认证，由于采用无密码格式，这里省略验证
	// 第三阶段请求信息
	// VER, CMD, RSV, ATYP, ADDR, PORT
	res, err := RequestAddressAuth(client, dstServer, socksClientAuthResponse)
	if err != nil {
		return
	}

	// 转发消息
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go func() {
		defer wg.Done()
		n, _ := dstServer.Write(*res)
		if n > 0 {
			utils.Logger.Info("发送成功")
		}
	}()

	go func() {
		defer wg.Done()
		io.Copy(client, dstServer)
	}()

	wg.Wait()
}
