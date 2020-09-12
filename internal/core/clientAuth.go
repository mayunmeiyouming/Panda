package core

import (
	"Panda/utils"
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
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

// RequestVersionAndMethodAuth 是第一阶段协议版本及认证方式
func RequestVersionAndMethodAuth(dstServer *net.TCPConn) (*SocksClientAuthResponse, error) {
	dstServer.Write([]byte{0x05, 0x01, 0x00})

	resp := make([]byte, 1024)
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

// RequestAddressAuth 第二阶段根据认证方式执行对应的认证
func RequestAddressAuth(client *net.TCPConn, dstServer *net.TCPConn, socksClientAuthResponse *SocksClientAuthResponse) (*[]byte, error) {
	buff := make([]byte, 1024)
	n, err := client.Read(buff)
	if err != nil {
		utils.Logger.Error(err)
		return nil, err
	}

	// 解析 http 地址
	re := bytes.NewReader(buff[:n])
	a := bufio.NewReader(re)
	request, _ := http.ReadRequest(a)
	utils.Logger.Info("请求地址: ", request.Host)
	dst := strings.Split(request.Host, ":")
	var dstAddr = dst[0]
	var dstPort = "80"
	if len(dst) == 2 && dst[1] != "" {
		dstPort = dst[1]
	}

	// 暂时只支持 HTTP
	if dstPort != "80" {
		return nil, errors.New("暂时只支持 HTTP")
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
		return nil, err
	}

	// 接收代理服务器返回的结果
	n, err = dstServer.Read(resp[0:])
	if err != nil {
		utils.Logger.Debug(dstServer.RemoteAddr(), err)
		return nil, err
	}
	if resp[0] != 0x05 || resp[1] != 0x00 {
		utils.Logger.Debug("第二阶段协商出错")
		return nil, err
	}
	utils.Logger.Debug("认证成功")

	res := buff[:n]
	return &res, nil
}
