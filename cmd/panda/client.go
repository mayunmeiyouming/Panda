package panda

import (
	"Panda/internal/config"
	"Panda/utils"
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

// Client 是 Panda 的 Client 模式的实际入口
func Client(localePort string, remotePort string) {
	config.InitConfiguration("config", "./configs/", &config.CONFIG)
	utils.InitLogger(config.CONFIG.LoggerConfig)

	// proxy server 地址
	serverAddr, err := net.ResolveTCPAddr("tcp", ":"+remotePort)
	if err != nil {
		utils.Log.Error(err)
	}
	utils.Log.Debug("连接远程服务器: ", remotePort+"....")

	// 监听本地
	listenAddr, err := net.ResolveTCPAddr("tcp", ":"+localePort)
	if err != nil {
		utils.Log.Error(err)
	}
	utils.Log.Debug("监听本地端口: ", localePort)

	listener, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		utils.Log.Error(err)
	}

	for {
		client, err := listener.AcceptTCP()
		if err != nil {
			utils.Log.Error(err)
		}
		go handleProxyRequest(client, serverAddr)
	}
}

func handleProxyRequest(client *net.TCPConn, serverAddr *net.TCPAddr) {

	// 远程连接IO
	dstServer, err := net.DialTCP("tcp", nil, serverAddr)
	if err != nil {
		utils.Log.Debug("远程服务器地址连接错误!!!")
		utils.Log.Debug(err)
		return
	}
	defer dstServer.Close()

	defer client.Close()

	// 和远程端建立安全信道
	wg := new(sync.WaitGroup)
	wg.Add(2)

	// socket5请求认证协商
	// 第一阶段协议版本及认证方式
	dstServer.Write([]byte{0x05, 0x01, 0x00})

	resp := make([]byte, 1024)
	n, err := dstServer.Read(resp)
	if err != nil {
		utils.Log.Error(err)
		return
	}
	if n == 0 {
		utils.Log.Error("协议错误,服务器返回为空")
		return
	}
	if resp[1] == 0x00 && n == 2 {
		utils.Log.Debug("第一阶段协商成功")
	} else {
		utils.Log.Error("协议错误，连接失败")
		return
	}
	// 第二阶段根据认证方式执行对应的认证，由于采用无密码格式，这里省略验证
	// 第三阶段请求信息
	// VER, CMD, RSV, ATYP, ADDR, PORT
	buff := make([]byte, 1024)
	n, err = client.Read(buff)
	if err != nil {
		utils.Log.Error(err)
		return
	}
	localReq := buff[:n]
	j := 0
	z := 0
	httpreq := []string{}
	for i := 0; i < n; i++ {
		if buff[i] == 32 {
			httpreq = append(httpreq, string(buff[j:i]))
			j = i + 1
		}
		if buff[i] == 10 {
			z += 1
		}
	}

	dstURI, err := url.ParseRequestURI(httpreq[1])
	if err != nil {
		utils.Log.Error(err)
		return
	}
	var dstAddr string
	var dstPort = "80"
	dstAddrPort := strings.Split(dstURI.Host, ":")
	if len(dstAddrPort) == 1 {
		dstAddr = dstAddrPort[0]
	} else if len(dstAddrPort) == 2 {
		dstAddr = dstAddrPort[0]
		dstPort = dstAddrPort[1]
	} else {
		utils.Log.Debug("URL parse error!")
		return
	}

	resp = []byte{0x05, 0x01, 0x00, 0x03}
	// 域名
	// dstAddrLenBuff := bytes.NewBuffer(make([]byte, 1))
	// binary.BigEndian.PutUint16(dstAddrLenBuff, uint8(len(dstAddr)))
	// binary.Write(dstAddrLenBuff, binary.BigEndian, uint8(len(dstAddr)))
	// log.Print("AdrrLength:", dstAddrLenBuff.Bytes()[dstAddrLenBuff.Len()-1])
	// resp = append(resp, dstAddrLenBuff.Bytes()[dstAddrLenBuff.Len()-1])
	resp = append(resp, byte(len([]byte(dstAddr))))
	resp = append(resp, []byte(dstAddr)...)
	// 端口
	dstPortBuff := bytes.NewBuffer(make([]byte, 0))
	dstPortInt, err := strconv.ParseUint(dstPort, 10, 16)
	if err != nil {
		utils.Log.Error(err)
		return
	}
	binary.Write(dstPortBuff, binary.BigEndian, dstPortInt)
	dstPortBytes := dstPortBuff.Bytes() // int为8字节
	resp = append(resp, dstPortBytes[len(dstPortBytes)-2:]...)
	n, err = dstServer.Write(resp)
	if err != nil {
		utils.Log.Debug(dstServer.RemoteAddr(), err)
		return
	}
	n, err = dstServer.Read(resp[0:])
	if err != nil {
		utils.Log.Debug(dstServer.RemoteAddr(), err)
		return
	}
	var targetResp [10]byte
	copy(targetResp[:10], resp[:n])
	specialResp := [10]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	if targetResp != specialResp {
		utils.Log.Debug("第二阶段协商出错")
		return
	}
	utils.Log.Debug("认证成功")

	// 转发消息
	go func() {
		defer wg.Done()
		dstServer.Write(localReq)
		// SecureCopy(localClient, dstServer, auth.Encrypt)
	}()

	go func() {
		defer wg.Done()
		io.Copy(client, dstServer)
	}()

	wg.Wait()

}
