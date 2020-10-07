package proxy

import (
	"Panda/pkg/socks"
	"Panda/utils"
	"net"
)

// TCPRemote is only server
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

		// utils.Logger.Debug("正在处理请求中")
		go func() {
			sc := shadow(client)

			target, err := socks.ReadAddr(sc)
			if err != nil {
				utils.Logger.Error("获取地址出错: ", err)
				return
			}

			utils.Logger.Info("代理地址: ", target.String())
			// 代理阶段
			destinationServer, err := net.Dial("tcp", target.String())
			if err != nil {
				return
			}
			defer destinationServer.Close()

			relay(destinationServer, sc)
		}()
	}
}

// TCPLocal is only client
// Listen on addr and proxy to server to reach target from getAddr.
func TCPLocal(addr, server string, shadow func(net.Conn) net.Conn, getAddr func(c net.Conn) (*socks.SocksAddressRequest, *byte, error)) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		utils.Logger.Info("failed to listen on ", addr, ": ", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			utils.Logger.Info("failed to accept: ", err)
			continue
		}

		go func() {
			defer c.Close()
			socksAddressRequest, _, err := getAddr(c)
			if err != nil {
				utils.Logger.Info("failed to get target address: ", err)
				return
			}

			rc, err := net.Dial("tcp", server)
			if err != nil {
				utils.Logger.Info("failed to connect to server ", server, ": ", err)
				return
			}
			defer rc.Close()

			rc = shadow(rc)

			addr := makeAddrRequest(socksAddressRequest.ADDR, socksAddressRequest.PORT)

			if _, err = rc.Write(addr); err != nil {
				utils.Logger.Info("failed to send target address: ", err)
				return
			}

			utils.Logger.Info("proxy ", c.RemoteAddr(), " <-> ", server, " <-> ", socksAddressRequest.ADDR)
			if err = relay(rc, c); err != nil {
				utils.Logger.Info("relay error: ", err)
			}
		}()
	}
}
