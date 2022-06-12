package main

import (
	"fmt"
	"http_server/src/handler"
	"http_server/src/router"
	"http_server/src/server_config"
	"http_server/src/service"
)

func main() {
	service.InitService()
	router.InitRouter()
	handler.InitHandler()
	service.GetHttpServer().Run(fmt.Sprintf("%v:%v", server_config.GlbserverConfig.HttpIp, server_config.GlbserverConfig.HttpPort))
}
