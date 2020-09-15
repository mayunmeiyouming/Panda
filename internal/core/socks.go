package core

import (
	"Panda/utils"
	"fmt"
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
		SecureCopy(destinationServer, conn)
		// io.Copy(destinationServer, conn)
	}()

	// 远程得到的内容copy到源地址
	go func() {
		defer wg.Done()
		SecureCopy(conn, destinationServer)
		// io.Copy(conn, destinationServer)
	}()
	wg.Wait()

	utils.Logger.Info("代理成功")
}

// SocksClient ...
func SocksClient(client *net.TCPConn, dstServer *net.TCPConn) {
	if client == nil || dstServer == nil {
		return
	}
	defer dstServer.Close()
	defer client.Close()

	// socket5请求认证协商
	res, port, err := SocksClientAuth(client, dstServer)
	if err != nil {
		return
	}

	// 转发消息
	if *port == 443 {
		fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		utils.Logger.Info("HTTP 2222包: ", string(*res))
		dstServer.Write(*res)
		utils.Logger.Debug("HTTP发送成功")
	} 
	
	//进行转发
	utils.Logger.Debug("数据转发中..........")
	// go SecureCopy(dstServer, client)
	// SecureCopy(client, dstServer)
	
	go io.Copy(dstServer, client)
	io.Copy(client, dstServer)
	

	utils.Logger.Info("代理成功")
}

// SecureCopy ...
func SecureCopy(dst io.ReadWriteCloser, src io.Reader) (written int64, err error) {
	size := 1024
	buf := make([]byte, size)

	for {
		utils.Logger.Debug("准备发送")
		nr, er := src.Read(buf[:])
		utils.Logger.Debug("发送长度: ", nr)
		// utils.Logger.Debug("代理数据: ", string(buf[:nr]))
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			utils.Logger.Debug("发送成功")
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
	}
	return written, err
}
