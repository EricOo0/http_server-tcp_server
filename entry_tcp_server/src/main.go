package main

import (
	"fmt"
	"net"
	"tcp_server/src/cache"
	_ "tcp_server/src/dao/impl"
	"tcp_server/src/db"
	"tcp_server/src/logger"
	"tcp_server/src/server_config"
	"tcp_server/src/service"
)

func main() {
	server_config.InitServerConfig()
	logger.BootLogger()
	cache.BootRedisClient()
	db.InitDB()
	//启动tcp服务器
	lis, err := net.Listen("tcp", fmt.Sprintf("%v:%v", server_config.GlbserverConfig.TcpHost, server_config.GlbserverConfig.TcpPort))
	if err != nil {
		logger.DefaultLogger.Error("init Tcpserver failed")
		panic(err)
	}
	defer func() {
		lis.Close()
	}()

	service.InitTcpServer()
	service.GetTcpServer().Serve(lis)
}
