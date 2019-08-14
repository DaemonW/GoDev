package api

import (
	"daemonw/controller"
	"daemonw/middleware"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func initRouter() {
	gin.SetMode(gin.ReleaseMode)
	router = newEngine()
	//init routers
	initUserRouter()
	//initStaticRouter()
}

func newEngine() *gin.Engine {
	engine := gin.New()
	engine.Use(middleware.Logger(), gin.Recovery())
	return engine
}

func GetRouter() *gin.Engine {
	initRouter()
	return router
}

func initUserRouter() {
	router.POST("/api/users", controller.CreateUser)
	router.GET("/api/user/:id", controller.GetUser)
	router.GET("/api/users", controller.GetUsers)
	router.POST("/api/user/token", controller.GenToken)
	router.PUT("/api/user/:id", controller.UpdateUser)

	authRouter := router.Group("")
	authRouter.Use(middleware.JwtAuth())
	authRouter.GET("/security/verify_code", controller.GetVerifyCode)
}
