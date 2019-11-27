package router

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
	plainRooter.Use(middleware.RateLimiter(entity.NewLimiter(*dao.Redis()), 1))
	plainRooter.POST("/api/users", controller.CreateUser)
	plainRooter.POST("/api/tokens", controller.GenToken)
	plainRooter.GET("/api/security/verify_codes", controller.GetVerifyCode)

	authRouter := r.Group("")
	authRouter.Use(middleware.JwtAuth())
	authRouter.Use(middleware.RateLimiter(entity.NewLimiter(*dao.Redis()), 1))
	authRouter.GET("/api/users/:id", controller.GetUser)
	authRouter.GET("/api/users", controller.GetUsers)
	authRouter.PUT("/api/users/:id", controller.UpdateUser)
	authRouter.DELETE("/api/users/:id", controller.DeleteUser)

	fileRouter := r.Group("")
	fileRouter.Use(middleware.RateLimiter(entity.NewLimiter(*dao.Redis()), 1))
	fileRouter.POST("/api/users/:id/files", controller.CreateFile)

	appRouter := r.Group("")
	appRouter.POST("/api/apps", controller.CreateApp)
	appRouter.POST("api/app", controller.QueryApp)
	appRouter.GET("/api/app/download", controller.DownloadApp)
}
