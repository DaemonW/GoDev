package api

import (
	"daemonw/controller"
	"daemonw/middleware"
	"github.com/gin-gonic/gin"
	"daemonw/conf"
	"time"
)

var router *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	router = newEngine()
	//init routers
	initUserRouter()
	initStaticRouter()
}

func newEngine() *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.Logger(), gin.Recovery())
	return engine
}

func GetRouter() *gin.Engine {
	return router
}

func initStaticRouter() {
	binDir := conf.BinDir
	router.StaticFile("/", binDir+"/static/html/index.html")
	router.Static("/static/", binDir+"/static")
	router.StaticFile("/user/", binDir+"/static/html/user.html")
}

func initUserRouter() {
	router.GET("/users/:user_id", controller.GetUser)
	router.POST("/users", controller.CreateUser)
	router.POST("/login", controller.Login)
	router.GET("/users", controller.GetAllUsers)

	authRouter := router.Group("")
	authRouter.Use(middleware.JwtAuth())
	authRouter.GET("/setting/limit", controller.LimitUserAccessCount)

	limitRouter := authRouter.Group("")
	limitRouter.Use(middleware.UserRateLimiter(2))
	limitRouter.Use(middleware.UserCountLimiter(20, time.Second*10))
	limitRouter.GET("/verify", controller.GetVerifyCode)
}
