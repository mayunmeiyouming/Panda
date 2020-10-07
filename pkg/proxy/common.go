package proxy

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

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

// 构造地址请求
func makeAddrRequest(host string, port string) []byte {
	addr := make([]byte, 0)

	address := net.ParseIP(host)
	if address != nil {
		// IPv4
		if len(address) == 4 {
			addr = append(addr, 0x01)
		} else {
			// IPv6
			addr = append(addr, 0x04)
		}
	} else {
		addr = append(addr, 0x03)
		addr = append(addr, byte(len([]byte(host)))) // 域名字节长度
	}

	// 域名
	addr = append(addr, []byte(host)...)

	// 端口
	b := []byte{0, 0}
	r, _ := strconv.Atoi(port)
	binary.BigEndian.PutUint16(b, uint16(r))
	addr = append(addr, b[:2]...)

	return addr
}
