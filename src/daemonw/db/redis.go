package db

import (
	"github.com/go-redis/redis"
	"daemonw/log"
	"daemonw/conf"
	"strconv"
)

var rc *redis.Client
var ring *redis.Ring

func init() {
	cfg := conf.Config.Redis
	addr := cfg.Host + ":" + strconv.Itoa(cfg.Port)
	clientOption:=&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.Num,
	}
	rc = redis.NewClient(clientOption)
	_, err := rc.Ping().Result()
	if err != nil {
		log.Fatal().Err(err).Msg("connect redis failed")
	}
	ringOptions:=&redis.RingOptions{
		Addrs: map[string]string{
			"server1": addr,
		},
	}
	ring = redis.NewRing(ringOptions)
}

func GetRedis() *redis.Client {
	return rc
}

func GetRing() *redis.Ring {
	return ring
}
