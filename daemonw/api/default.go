package api

import (
	"daemonw/controller"
	"daemonw/dao"
	"daemonw/entity"
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
	router.POST("/api/users",
		middleware.RateLimiter(entity.NewLimiter(*dao.Redis()), 1),
		controller.CreateUser)
	router.POST("/api/user/token",
		middleware.RateLimiter(entity.NewLimiter(*dao.Redis()), 1),
		controller.GenToken)
	router.GET("/api/security/verify_code",
		middleware.RateLimiter(entity.NewLimiter(*dao.Redis()), 1),
		controller.GetVerifyCode)

	authRouter := router.Group("")
	authRouter.Use(middleware.JwtAuth())
	authRouter.GET("/api/user/:id",
		middleware.RateLimiter(entity.NewLimiter(*dao.Redis()), 1),
		controller.GetUser)
	authRouter.GET("/api/users",
		middleware.RateLimiter(entity.NewLimiter(*dao.Redis()), 1),
		controller.GetUsers)
	authRouter.PUT("/api/user/:id",
		middleware.RateLimiter(entity.NewLimiter(*dao.Redis()), 1),
		controller.UpdateUser)
	authRouter.DELETE("/api/user/:id",
		middleware.RateLimiter(entity.NewLimiter(*dao.Redis()), 1),
		controller.DeleteUser)
}
