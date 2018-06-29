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
	router = gin.Default()

	//init routers
	initUserRouter()
	initPageRouter()
}

func GetRouter() *gin.Engine {
	return router
}

func initPageRouter() {
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
	authRouter := router.Group("api")
	authRouter.Use(middleware.JwtAuth())
	authRouter.Use(middleware.UserRateLimiter(2))
	authRouter.Use(middleware.UserCountLimiter(20, time.Second*10))
	authRouter.GET("/verify", controller.GetVerifyCode)
}
