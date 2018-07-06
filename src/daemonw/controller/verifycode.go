package controller

import (
	"github.com/gin-gonic/gin"
	"daemonw/util"
	"daemonw/db"
	"strconv"
	"time"
	"net/http"
	"daemonw/model"
	"daemonw/log"
)

func GetVerifyCode(c *gin.Context) {
	uid := c.MustGet("uid").(uint64)
	code := util.RandomNum(8)
	err := db.GetRedis().Set("verify_code:"+strconv.FormatUint(uid, 10), code, time.Minute*10).Err()
	if err != nil {
		log.Error().Err(err).Msg("generate verify code failed")
		c.JSON(http.StatusInternalServerError, model.NewResp().SetError(model.ErrInternalServer))
		return
	}
	c.JSON(http.StatusOK, model.NewResp().AddResult("verify_code", code))
}
