package panda

import (
	"Panda/internal/config"
	"Panda/internal/core"
	"Panda/utils"
	"net"
)

// Server 是 Panda 的 Server 模式的实际入口
func Server(port string) {
	config.InitConfiguration("config", "./configs/", &config.CONFIG)
	utils.InitLogger(config.CONFIG.LoggerConfig)

	listenAddr, err := net.ResolveTCPAddr("tcp", ":"+port)
	l, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		utils.Log.Error("监听端口失败，端口可能被占用")
	}

	for {
		utils.Log.Debug("等待连接")
		client, err := l.AcceptTCP()
		if err != nil {
			utils.Log.Error(err)
		}
		utils.Log.Debug("正在处理请求中")
		go handleClientRequest(client)
	}
}

func handleClientRequest(client *net.TCPConn) {
	if client == nil {
		return
	}
	defer client.Close()
	core.SocksServe(client)
}
