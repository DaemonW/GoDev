package dao

import (
	"daemonw/conf"
	"daemonw/entity"
	"daemonw/util"
	"daemonw/xlog"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"strconv"
)

var (
	UserDao *userDao
	FileDao *fileDao
)

const (
	DialWithoutPass = "postgres://%s@%s:%d/%s?sslmode=%s"
	DialWithPass    = "postgres://%s:%s@%s:%d/%s?sslmode=%s"
)

var (
	dbConn *sqlx.DB
	rsConn *redis.Client
)

func initDB() {
	var err error
	c := &conf.Config.Database
	//connStr := "postgres://postgres:a123456@localhost:5432/mydb?sslmode=disable"
	var connParams string
	if c.Password == "" {
		connParams = fmt.Sprintf(DialWithoutPass, c.User, c.Host, c.Port, c.Name, c.SSLMode)
	} else {
		connParams = fmt.Sprintf(DialWithPass, c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
	}
	dbConn, err = sqlx.Connect("postgres", connParams)
	util.PanicIfErr(err)
}

func initUser() {
	admin:= entity.NewUser("admin","admin")
	admin.Status = entity.UserStatusNormal
	admin.Role = entity.UserRoleAdmin
	err:= UserDao.CreateUserIfNotExist(admin)
	util.PanicIfErr(err)
}

func initRedis() {
	cfg := conf.Config.Redis
	addr := cfg.Host + ":" + strconv.Itoa(cfg.Port)
	clientOption := &redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.Index,
	}
	rsConn = redis.NewClient(clientOption)
	_, err := rsConn.Ping().Result()
	util.PanicIfErr(err)
}

func DB() *sqlx.DB {
	return dbConn;
}

func Redis() *redis.Client {
	return rsConn
}

func InitDao() {
	initDB()
	initRedis()
	UserDao = newUserDao()
	FileDao = newFileDao()
	initUser()
}

func CloseDao() {
	if dbConn != nil {
		xlog.Info().Msg("close database connection")
		util.PanicIfErr(dbConn.Close())
	}
	if rsConn != nil {
		xlog.Info().Msg("close redis connection")
		util.PanicIfErr(rsConn.Close())
	}
}
