package dao

import (
	"daemonw/conf"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"strconv"
)

var UserDao *userDao

func InitDaoManager() {
	UserDao = newUserDao()
}

const (
	DialWithoutPass = "postgres://%s@%s:%d/%s?sslmode=%s"
	DialWithPass    = "postgres://%s:%s@%s:%d/%s?sslmode=%s"
)

var (
	dbConn *sqlx.DB
	rsConn *redis.Client
)

func InitDB() error {
	var err error
	c := &conf.Config.Database
	//connStr := "postgres://postgres:a123456@localhost:5432/mydb?sslmode=disable"
	var connParams string
	if c.Password == "" {
		connParams = fmt.Sprintf(DialWithoutPass, c.User, c.Host, c.Port, c.Name, c.SSLMode)
	} else {
		connParams = fmt.Sprintf(DialWithPass, c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
	}
	sqlx.Rebind(sqlx.DOLLAR, "?")
	dbConn, err = sqlx.Connect("postgres", connParams)
	return err
}

func InitRedis() error {
	cfg := conf.Config.Redis
	addr := cfg.Host + ":" + strconv.Itoa(cfg.Port)
	clientOption := &redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.Num,
	}
	rsConn = redis.NewClient(clientOption)
	_, err := rsConn.Ping().Result()
	return err
}

func DB() *sqlx.DB {
	return dbConn;
}

func Redis() *redis.Client {
	return rsConn
}
