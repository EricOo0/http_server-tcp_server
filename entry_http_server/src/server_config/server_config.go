package server_config

import (
	"github.com/go-ini/ini"
	constant "http_server/src/common"
	"os"
	"path/filepath"
)

var GlbserverConfig ServerConfig

type ServerConfig struct {
	HttpIp   string
	HttpPort string

	RedisHost string
	RedisPort string
	RedisDB   int
	RedisPwd  string

	TcpHost string
	TcpPort string
}

func ConfigFile() string {
	tempDirPath, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}
	for {
		//已经找到根目录
		if tempDirPath == "/" {
			panic("can not find project base path")
		}
		configPath := filepath.Join(tempDirPath, "config")
		_, err1 := os.Stat(configPath)
		if err1 == nil {
			configFilePath := filepath.Join(configPath, constant.ConfigFileName)
			return configFilePath
		}
		//不存在往上级目录查找
		tempDirPath = filepath.Dir(tempDirPath)
	}
	panic("can not find project bath path")
}
func LoadConfigFile(configFilePath string) *ini.File {
	fileAbsPath, err := filepath.Abs(configFilePath)
	if err != nil {
		panic(err)
	}
	file, err := ini.Load(fileAbsPath)
	if err != nil {
		panic(err)
	}
	return file
}
func InitServerConfig() {
	configFilepath := ConfigFile()
	configFile := LoadConfigFile(configFilepath)

	section := configFile.Section("http")
	GlbserverConfig.HttpIp = section.Key("ip").String()
	GlbserverConfig.HttpPort = section.Key("port").String()

	section = configFile.Section("redis")
	GlbserverConfig.RedisHost = section.Key("host").String()
	GlbserverConfig.RedisPort = section.Key("port").String()
	GlbserverConfig.RedisPwd = section.Key("password").String()
	GlbserverConfig.RedisDB, _ = section.Key("DB").Int()

	section = configFile.Section("tcp_server")
	GlbserverConfig.TcpHost = section.Key("ip").String()
	GlbserverConfig.TcpPort = section.Key("port").String()

}
