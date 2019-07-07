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
	router.GET("/api/user/:user_id", controller.GetUser)
	router.POST("/api/users", controller.CreateUser)
	router.POST("/api/user/auth/token", controller.Login)
	router.POST("/api/user/auth/active", controller.ActiveUser)
	router.PUT("/api/user/:user_id", controller.ActiveUser)
	router.GET("/api/users", controller.GetUsers)

	authRouter := router.Group("")
	authRouter.Use(middleware.JwtAuth())
	authRouter.GET("/security/verify_code", controller.GetVerifyCode)
}
