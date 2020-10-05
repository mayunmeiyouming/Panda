package socks

import (
	"Panda/utils"
	"bufio"
	"bytes"
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// ParseHTTP ...
func ParseHTTP(client net.Conn) ([]byte, *string, *string, error) {
	buff := make([]byte, 0, 2048)
	reader := bufio.NewReader(client)
	for {
		temp, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, nil, err
		}

		buff = append(buff[:len(buff)], temp[:len(temp)]...)

		if len(temp) == 2 {
			break
		}
	}

	utils.Logger.Info("HTTP包: ", string(buff[:len(buff)]))

	// 解析 http 地址
	re := bytes.NewReader(buff[:len(buff)])
	a := bufio.NewReader(re)
	request, err := http.ReadRequest(a)
	if err != nil {
		utils.Logger.Info("HTTP 包有错误: ", string(buff))
		return nil, nil, nil, err
	}

	utils.Logger.Info("请求地址: ", request.Host)
	dst := strings.Split(request.Host, ":")
	var dstAddr = dst[0]
	var dstPort = "80"
	if len(dst) == 2 && dst[1] != "" {
		dstPort = dst[1]
	}

	return buff, &dstAddr, &dstPort, nil
}

// SOCKS address types as defined in RFC 1928 section 5.
const (
	AtypIPv4       = 1
	AtypDomainName = 3
	AtypIPv6       = 4
)

// MaxAddrLen is the maximum size of SOCKS address in bytes.
const MaxAddrLen = 1 + 1 + 255 + 2

// Addr represents a SOCKS address as defined in RFC 1928 section 5.
type Addr []byte

// String serializes SOCKS address a to string form.
func (a Addr) String() string {
	var host, port string

	switch a[0] { // address type
	case AtypDomainName:
		host = string(a[2 : 2+int(a[1])])
		port = strconv.Itoa((int(a[2+int(a[1])]) << 8) | int(a[2+int(a[1])+1]))
	case AtypIPv4:
		host = net.IP(a[1 : 1+net.IPv4len]).String()
		port = strconv.Itoa((int(a[1+net.IPv4len]) << 8) | int(a[1+net.IPv4len+1]))
	case AtypIPv6:
		host = net.IP(a[1 : 1+net.IPv6len]).String()
		port = strconv.Itoa((int(a[1+net.IPv6len]) << 8) | int(a[1+net.IPv6len+1]))
	}

	return net.JoinHostPort(host, port)
}

func readAddr(r io.Reader, b []byte) (Addr, error) {
	if len(b) < MaxAddrLen {
		return nil, io.ErrShortBuffer
	}
	_, err := io.ReadFull(r, b[:1]) // read 1st byte for address type
	if err != nil {
		return nil, err
	}

	switch b[0] {
	case AtypDomainName:
		_, err = io.ReadFull(r, b[1:2]) // read 2nd byte for domain length
		if err != nil {
			return nil, err
		}
		_, err = io.ReadFull(r, b[2:2+int(b[1])+2])
		return Addr(b[:1+1+int(b[1])+2]), err
	case AtypIPv4:
		_, err = io.ReadFull(r, b[1:1+net.IPv4len+2])
		return Addr(b[:1+net.IPv4len+2]), err
	case AtypIPv6:
		_, err = io.ReadFull(r, b[1:1+net.IPv6len+2])
		return Addr(b[:1+net.IPv6len+2]), err
	}

	return nil, errors.New("address not support")
}

// ReadAddr reads just enough bytes from r to get a valid Addr.
func ReadAddr(r io.Reader) (Addr, error) {
	return readAddr(r, make(Addr, MaxAddrLen))
}
