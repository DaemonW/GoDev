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
	router.POST("/api/user/token", controller.GenToken)
	router.GET("/api/security/verify_code", controller.GetVerifyCode)

	authRouter := router.Group("")
	authRouter.Use(middleware.JwtAuth())
	authRouter.GET("/api/user/:id", controller.GetUser)
	authRouter.GET("/api/users", controller.GetUsers)
	authRouter.PUT("/api/user/:id", controller.UpdateUser)
	authRouter.DELETE("/api/user/:id", controller.DeleteUser)
}
