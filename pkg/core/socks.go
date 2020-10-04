package core

import (
	"Panda/utils"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

// SocksServe ...
func SocksServe(conn net.Conn) {
	// 协商阶段
	socksAddressRequest, _, err := SocksAuth(conn)
	if err != nil {
		return
	}

	// 代理阶段
	destinationServer, err := net.DialTCP("tcp", nil, socksAddressRequest.RAWADDR)
	if err != nil {
		return
	}
	defer destinationServer.Close()

	relay(destinationServer, conn)

	utils.Logger.Info("代理成功")
}

// SocksClient ...
func SocksClient(client *net.TCPConn, dstServer net.Conn) {

	// socket5请求认证协商
	res, port, err := SocksClientAuth(client, dstServer)
	if err != nil {
		return
	}

	// dstCrypto := getCrypt(method)

	// 转发消息
	if *port == 443 {
		fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		dstServer.Write(*res)
		utils.Logger.Debug("HTTP发送成功")
	}

	//进行转发
	utils.Logger.Debug("数据转发中..........")

	relay(dstServer, client)

	utils.Logger.Info("代理成功")
}

// relay copies between left and right bidirectionally
func relay(left, right net.Conn) error {
	var err, err1 error
	var wg sync.WaitGroup
	var wait = 5 * time.Second
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err1 = io.Copy(right, left)
		right.SetReadDeadline(time.Now().Add(wait)) // unblock read on right
	}()
	_, err = io.Copy(left, right)
	left.SetReadDeadline(time.Now().Add(wait)) // unblock read on left
	wg.Wait()
	if err1 != nil && !errors.Is(err1, os.ErrDeadlineExceeded) { // requires Go 1.15+
		return err1
	}
	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) {
		return err
	}
	return nil
}
