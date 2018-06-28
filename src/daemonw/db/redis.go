package db

import (
	"github.com/go-redis/redis"
	"daemonw/log"
	"daemonw/conf"
	"strconv"
)

var rc *redis.Client

func init() {
	cfg := conf.Config.Redis
	addr := cfg.Host + ":" + strconv.Itoa(cfg.Port)
	rc = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.Num,
	})

	_, err := rc.Ping().Result()
	if err != nil {
		log.Fatal().Err(err).Msg("connect redis failed")
	}
}

func GetRedis() *redis.Client {
	return rc
}
