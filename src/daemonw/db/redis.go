package db

import (
	"daemonw/conf"
	"github.com/go-redis/redis"
	"strconv"
)

var rc *redis.Client

func InitRedis() error {
	cfg := conf.Config.Redis
	addr := cfg.Host + ":" + strconv.Itoa(cfg.Port)
	clientOption := &redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.Num,
	}
	rc = redis.NewClient(clientOption)
	_, err := rc.Ping().Result()
	return err
}

func GetRedis() *redis.Client {
	return rc
}
