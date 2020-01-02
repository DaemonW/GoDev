package router

import (
	"daemonw/conf"
	"daemonw/controller"
	"daemonw/dao"
	"daemonw/entity"
	"daemonw/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
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
	engine.Use(middleware.AllowCORS())
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
	appAdminRouter.Use(middleware.JwtAuth())
	appAdminRouter.POST("/api/admin/apps", controller.CreateApp)
	appAdminRouter.DELETE("/api/admin/app/:id", controller.DeleteApp)
	appAdminRouter.PUT("api/admin/app/:id", controller.UpdateApp)

	appRouter := r.Group("")
	appRouter.GET("/api/apps", controller.QueryApps)
	appRouter.GET("/api/download/app/:id", controller.DownloadApp)
	appRouter.GET("/api/app/detail/:id", controller.GetAppInfo)
	appRouter.Static("/api/app/resources", conf.Config.Data+"/res")
	appRouter.Static("/static", conf.Config.Static)
	appRouter.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/static/index.html")
	})
}
