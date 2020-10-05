package proxy

import (
	"Panda/pkg/socks"
	"Panda/utils"
	"net"
)

// TCPRemote ...
func TCPRemote(addr string, shadow func(net.Conn) net.Conn) {
	listenAddr, err := net.ResolveTCPAddr("tcp", addr)
	l, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		utils.Logger.Fatal("监听端口失败，端口可能被占用")
	}

	for {
		utils.Logger.Info("等待连接")
		client, err := l.AcceptTCP()
		if err != nil {
			utils.Logger.Error(err)
		}
		defer client.Close()

		utils.Logger.Debug("正在处理请求中")
		go func() {
			sc := shadow(client)

			// // 协商阶段
			// socksAddressRequest, _, err := socks.SocksAuth(sc)
			// if err != nil {
			// 	return
			// }
			tgt, err := socks.ReadAddr(sc)

			utils.Logger.Info("代理地址: ", string(addr))
			// 代理阶段
			destinationServer, err := net.Dial("tcp", tgt.String())
			if err != nil {
				return
			}
			defer destinationServer.Close()

			relay(destinationServer, sc)

			utils.Logger.Info("代理成功")
		}()
	}
}

// TCPLocal ...
// Listen on addr and proxy to server to reach target from getAddr.
func TCPLocal(addr, server string, shadow func(net.Conn) net.Conn, getAddr func(c net.Conn) (*socks.SocksAddressRequest, *byte, error)) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		utils.Logger.Info("failed to listen on %s: %v", addr, err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			utils.Logger.Info("failed to accept: %s", err)
			continue
		}

		go func() {
			defer c.Close()
			socksAddressRequest, _, err := getAddr(c)
			if err != nil {
				utils.Logger.Info("failed to get target address: %v", err)
				return
			}

			rc, err := net.Dial("tcp", server)
			if err != nil {
				utils.Logger.Info("failed to connect to server %v: %v", server, err)
				return
			}
			defer rc.Close()

			rc = shadow(rc)

			if _, err = rc.Write(socksAddressRequest.DSTADDR); err != nil {
				utils.Logger.Info("failed to send target address: %v", err)
				return
			}

			utils.Logger.Info("proxy %s <-> %s <-> %s", c.RemoteAddr(), server, socksAddressRequest.DSTADDR)
			// if err = relay(rc, c); err != nil {
			// 	utils.Logger.Info("relay error: %v", err)
			// }
		}()
	}
}
