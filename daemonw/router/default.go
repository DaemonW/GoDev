package router

import (
	"daemonw/conf"
	"daemonw/controller"
	"daemonw/dao"
	"daemonw/entity"
	"daemonw/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"path/filepath"
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
	authRouter.GET("/api/user/:id", controller.GetUser)
	authRouter.GET("/api/users", controller.GetUsers)
	authRouter.PUT("/api/user/:id", controller.UpdateUser)
	authRouter.DELETE("/api/user/:id", controller.DeleteUser)

	fileRouter := r.Group("")
	fileRouter.Use(middleware.RateLimiter(entity.NewLimiter(*dao.Redis()), 1))
	fileRouter.POST("/api/user/:id/files", controller.CreateFile)

	appAdminRouter := r.Group("")
	appAdminRouter.POST("/api/admin/apps", controller.CreateApp)

	appRouter := r.Group("")
	appRouter.GET("/api/apps", controller.QueryApps)
	appRouter.GET("/api/app/:id/downloads", controller.DownloadApp)
	fmt.Println(filepath.Dir(conf.Config.Data)+"/web")
	appRouter.Static("/api/static", filepath.Dir(conf.Config.Data)+"/web")

	resRouter := r.Group("")
	resRouter.Use(middleware.ResourceAuth())
	resRouter.Static("/api/resource/app/downloads", conf.Config.Data)
}
