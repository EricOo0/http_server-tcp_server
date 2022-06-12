package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"http_server/src/cache"
	"http_server/src/logger"
	api "http_server/src/proto/proto"
	"http_server/src/server_config"
)

var httpEngine *gin.Engine
var tcpClient api.AdminServiceClient

func GetHttpServer() *gin.Engine {
	return httpEngine
}
func GetTcpClient() api.AdminServiceClient {
	return tcpClient
}
func InitService() {
	server_config.InitServerConfig()
	logger.BootLogger()
	cache.BootRedisClient()
	httpEngine = gin.Default()

	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", server_config.GlbserverConfig.TcpHost, server_config.GlbserverConfig.TcpPort), grpc.WithInsecure())
	if err != nil {
		logger.DefaultLogger.Error("Dail Tcp server failed!")
		return
	}
	TcpClientConn := conn
	tcpClient = api.NewAdminServiceClient(TcpClientConn)
	logger.DefaultLogger.Info(fmt.Sprintf("Tcp client listen on %v:%v", server_config.GlbserverConfig.TcpHost, server_config.GlbserverConfig.TcpPort))

}
