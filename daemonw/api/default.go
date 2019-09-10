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
	initUserRouter(router)
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

func initUserRouter(r *gin.Engine) {
	plainRooter := r.Group("")
	plainRooter.POST("/api/users",
		middleware.RateLimiter(entity.NewLimiter(*dao.Redis()), 1),
		controller.CreateUser)
	plainRooter.POST("/api/tokens",
		middleware.RateLimiter(entity.NewLimiter(*dao.Redis()), 1),
		controller.GenToken)
	plainRooter.GET("/api/security/verify_codes",
		middleware.RateLimiter(entity.NewLimiter(*dao.Redis()), 1),
		controller.GetVerifyCode)

	authRouter := r.Group("")
	authRouter.Use(middleware.JwtAuth())
	authRouter.GET("/api/users/:id",
		middleware.RateLimiter(entity.NewLimiter(*dao.Redis()), 1),
		controller.GetUser)
	authRouter.GET("/api/users",
		middleware.RateLimiter(entity.NewLimiter(*dao.Redis()), 1),
		controller.GetUsers)
	authRouter.PUT("/api/users/:id",
		middleware.RateLimiter(entity.NewLimiter(*dao.Redis()), 1),
		controller.UpdateUser)
	authRouter.DELETE("/api/users/:id",
		middleware.RateLimiter(entity.NewLimiter(*dao.Redis()), 1),
		controller.DeleteUser)


	fileRouter := r.Group("")
	fileRouter.POST("/api/users/:id/files",
		middleware.RateLimiter(entity.NewLimiter(*dao.Redis()), 1),
		controller.CreateFile)
}
