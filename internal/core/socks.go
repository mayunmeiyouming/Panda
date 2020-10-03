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

	dstCrypto := NoCrpt{}
	clientCrypto := NoCrpt{}

	// 本地的内容copy到远程端
	go func() {
		defer wg.Done()
		SecureCopy(destinationServer, dstCrypto, conn, clientCrypto)
	}()

	// 远程得到的内容copy到源地址
	go func() {
		defer wg.Done()
		SecureCopy(conn, clientCrypto, destinationServer, dstCrypto)
	}()
	wg.Wait()

	utils.Logger.Info("代理成功")
}

// SocksClient ...
func SocksClient(client *net.TCPConn, dstServer *net.TCPConn, method byte) {
	if client == nil || dstServer == nil {
		return
	}
	defer dstServer.Close()
	defer client.Close()

	// socket5请求认证协商
	res, port, err := SocksClientAuth(client, dstServer, method)
	if err != nil {
		return
	}

	dstCrypto := NoCrpt{}
	clientCrypto := NoCrpt{}

	// 转发消息
	if *port == 443 {
		fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		dstServer.Write(dstCrypto.Encrypt(*res))
		utils.Logger.Debug("HTTP发送成功")
	} 
	
	//进行转发
	utils.Logger.Debug("数据转发中..........")
	
	go SecureCopy(dstServer, dstCrypto, client, clientCrypto)
	SecureCopy(client, clientCrypto, dstServer, dstCrypto)

	utils.Logger.Info("代理成功")
}

// SecureCopy ...
func SecureCopy(dst io.ReadWriter, dstCrypto Crypto, src io.ReadWriter, srcCrypto Crypto) (written int64, err error) {
	size := 20480
	buf := make([]byte, size)

	for {
		// utils.Logger.Debug("准备发送")
		nr, er := srcCrypto.DecodeRead(src, buf)
		// nr, er := src.Read(buf)
		// utils.Logger.Debug("发送长度: ", nr)
		// utils.Logger.Debug("代理数据: ", string(buf[:nr]))
		if nr > 0 {
			nw, ew := dstCrypto.EncodeWrite(dst, buf[0:nr])
			// nw, ew := dst.Write(buf[0:nr])
			// utils.Logger.Debug("发送成功")

			// 动态拓展切片长度
			if nr * 2 >= size && size <= 408600 {
				size = size * 2
				buf = make([]byte, size)
			}

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
