package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"

	"tcp_server/src/logger"
	"tcp_server/src/server_config"
)

var Idb *xorm.Engine

func GetDbInstance() *xorm.Engine {
	return Idb
}
func InitDB() {
	dbInstance, err := xorm.NewEngine(server_config.GlbserverConfig.DbDriver, server_config.GlbserverConfig.DbSource)
	if err != nil {
		logger.DefaultLogger.Error("Init db error")
	}
	dbInstance.SetMaxOpenConns(100)
	dbInstance.SetMaxIdleConns(10)
	dbInstance.ShowSQL(true)
	Idb = dbInstance
	logger.DefaultLogger.Info("Init db successfully")

}
