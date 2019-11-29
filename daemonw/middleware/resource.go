package middleware

import (
	"daemonw/dao"
	"daemonw/entity"
	"daemonw/xerr"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ResourceAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := c.Query("uuid")
		if uuid == "" {
			c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeQueryApp, "illegal request"))
			c.Abort()
			return
		}
		verifyCode := c.Query("c")
		if verifyCode == "" {
			c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeDownloadApp, "illegal request"))
			c.Abort()
			return
		}
		code := dao.Redis().Get("app:" + uuid).Val()
		if verifyCode != code {
			c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeDownloadApp, "illegal request"))
			c.Abort()
			return
		}
	}
}
