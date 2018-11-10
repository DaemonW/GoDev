package api

import (
	"daemonw/conf"
	"daemonw/controller"
	"daemonw/middleware"
	"github.com/gin-gonic/gin"
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
	authRouter.GET("/setting/limit")
	authRouter.GET("/verify", controller.GetVerifyCode)
}
