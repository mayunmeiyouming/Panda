package proxy

import (
	"Panda/utils"
	"net"

	"Panda/pkg/socks"
)

// SocksLocal ...
// Create a SOCKS server listening on addr and proxy to server.
func SocksLocal(addr, server string, shadow func(net.Conn) net.Conn) {
	utils.Logger.Info("SOCKS proxy ", addr, " <-> ", server)
	TCPLocal(addr, server, shadow, func(c net.Conn) (*socks.SocksAddressRequest, *byte, error) { return socks.SocksAuth(c) })
}
