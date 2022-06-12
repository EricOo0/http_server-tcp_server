package cache

import (
	"github.com/go-redis/redis"
	"tcp_server/src/logger"
	"time"
)

// 声明一个全局的rdb变量
var Rdb *redis.Client

// 初始化连接
func BootRedisClient() (err error) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		Password:     "", // no password set
		DB:           1,  // use default DB
		DialTimeout:  4 * (time.Second),
		ReadTimeout:  2 * (time.Second),
		WriteTimeout: 2 * (time.Second),

		//连接池容量及闲置连接数量
		PoolSize:     100, // 连接池最大socket连接数，默认为4倍CPU数， 4 * runtime.NumCPU
		MinIdleConns: 25,  //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；

	})

	_, err = Rdb.Ping().Result()
	if err != nil {
		logger.DefaultLogger.Info("Boot Redis client failed")
		panic(err)
	}
	logger.DefaultLogger.Info("Boot Redis client successfully")

	return nil
}
