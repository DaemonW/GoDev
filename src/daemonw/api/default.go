package api

import (
	"daemonw/controller"
	"daemonw/middleware"
	"github.com/gin-gonic/gin"
	"daemonw/conf"
)

var router *gin.Engine

func init() {

	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
	router.Use(middleware.ApiCounter)

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
	lp := middleware.NewLimiterPool()
	router.GET("/users/:userId", controller.GetUser)
	router.POST("/users", controller.CreateUser)
	router.POST("/login", controller.Login)
	router.GET("/users", middleware.IpLimiter(lp, 1, 2), controller.GetAllUsers)
	authRouter := router.Group("api")
	authRouter.Use(middleware.JwtAuth())
	authRouter.GET("/verify", controller.GetVerifyCode)
}
