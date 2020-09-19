package core

import (
	"Panda/utils"
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
)

// SocksClientAuthResponse ...
type SocksClientAuthResponse struct {
	VER    uint8
	METHOD uint8
}

// SocksClientAuth 是协商认证阶段
func SocksClientAuth(client *net.TCPConn, dstServer *net.TCPConn) (*[]byte, *int, error) {
	// 第一阶段协议版本及认证方式
	socksClientAuthResponse, err := RequestVersionAndMethodAuth(dstServer)
	if err != nil {
		return nil, nil, err
	}

	// 第二阶段根据认证方式执行对应的认证，由于采用无密码格式，这里省略验证，返回第三阶段请求信息
	// VER, CMD, RSV, ATYP, ADDR, PORT
	res, port, err := RequestAddressAuth(client, dstServer, socksClientAuthResponse)
	if err != nil {
		utils.Logger.Error("第二阶段出错: ", err)
		return nil, nil, err
	}

	return res, port, nil
}

// RequestVersionAndMethodAuth 是第一阶段协议版本及认证方式
func RequestVersionAndMethodAuth(dstServer *net.TCPConn) (*SocksClientAuthResponse, error) {
	dstServer.Write([]byte{0x05, 0x01, 0x00})

	resp := make([]byte, 2)
	n, err := dstServer.Read(resp)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}
	if n != 2 {
		utils.Logger.Error("协议错误,服务器返回为空")
		return nil, err
	}
	socksClientAuthResponse := SocksClientAuthResponse{}
	if resp[0] == 0x05 {
		utils.Logger.Debug("第一阶段协商成功")

		socksClientAuthResponse = SocksClientAuthResponse{
			VER:    uint8(resp[0]),
			METHOD: uint8(resp[1]),
		}
	} else {
		utils.Logger.Error("协议错误，连接失败")
		return nil, err
	}

	return &socksClientAuthResponse, nil
}

// RequestAddressAuth 第二阶段根据认证方式执行对应的认证，返回第三阶段请求信息
func RequestAddressAuth(client *net.TCPConn, dstServer *net.TCPConn, socksClientAuthResponse *SocksClientAuthResponse) (*[]byte, *int, error) {
	buff := make([]byte, 0)

	n := 0
	reader := bufio.NewReader(client)
	for {
		temp, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, err
		}
		buff = append(buff[:n], temp[:len(temp)]...)
		n += len(temp)

		if len(temp) == 2 {
			break
		}
	}

	utils.Logger.Info("HTTP包: ", string(buff[:n]))

	// 解析 http 地址
	re := bytes.NewReader(buff[:n])
	a := bufio.NewReader(re)
	request, err := http.ReadRequest(a)
	if err != nil {
		utils.Logger.Info("HTTP 包有错误: ", string(buff))
		return nil, nil, err
	}
	utils.Logger.Info("请求地址: ", request.Host)
	dst := strings.Split(request.Host, ":")
	var dstAddr = dst[0]
	var dstPort = "80"
	if len(dst) == 2 && dst[1] != "" {
		dstPort = dst[1]
	}

	// 暂时默认域名
	resp := []byte{0x05, 0x01, 0x00, 0x03}
	resp = append(resp, byte(len([]byte(dstAddr)))) // 域名字节长度
	resp = append(resp, []byte(dstAddr)...)         // 域名
	// 端口
	b := []byte{0, 0}
	r, _ := strconv.Atoi(dstPort)
	utils.Logger.Info("端口: ", r)
	utils.Logger.Info("Domain: ", dstAddr)
	binary.BigEndian.PutUint16(b, uint16(r))
	resp = append(resp, b[:2]...)

	// 发送地址给 Proxy Server
	n, err = dstServer.Write(resp)
	if err != nil {
		utils.Logger.Debug(dstServer.RemoteAddr(), err)
		return nil, nil, err
	}

	// 接收代理服务器返回的结果
	_, err = parseAddressResponse(dstServer)
	if err != nil {
		return nil, nil, err
	}

	return &buff, &r, nil
}

func parseAddressResponse(server *net.TCPConn) (*int, error) {
	resp := make([]byte, 2048)
	n, err := server.Read(resp[:])
	if err != nil {
		utils.Logger.Debug(server.RemoteAddr(), err)
		return nil, err
	}

	// 协议版本
	if resp[0] != 0x05 {
		utils.Logger.Debug("第二阶段协商出错")
		return nil, err
	}

	// 状态码
	if resp[1] != 0x00 {
		errMsg := "第二阶段协商出错，状态码: " + string(resp[1])
		utils.Logger.Debug(errMsg)
		return nil, errors.New(errMsg)
	}

	utils.Logger.Debug("第二阶段协商成功")
	utils.Logger.Debug("认证成功")

	return &n, nil
}
