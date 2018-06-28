package controller

import (
	"github.com/gin-gonic/gin"
	"daemonw/util"
	"daemonw/db"
	"strconv"
	"time"
	"net/http"
	"daemonw/model"
)

func GetVerifyCode(c *gin.Context) {
	uid := c.MustGet("uid").(uint64)
	code := util.RandomNum(8)
	err := db.GetRedis().Set("verify_cdoe-"+strconv.FormatUint(uid, 10), code, time.Minute*10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewResp().SetError(model.ErrInternalServer))
		return
	}
	c.JSON(http.StatusOK, model.NewResp().AddResult("verfy_code", code))
}
