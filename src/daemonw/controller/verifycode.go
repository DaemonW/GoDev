package controller

import (
	myerr "daemonw/errors"
	"github.com/gin-gonic/gin"
	"daemonw/util"
	"daemonw/db"
	"time"
	"net/http"
	"daemonw/model"
	"daemonw/log"
	"strings"
	"fmt"
)

var (
	SUPPORT_VERIFY = []string{"ACTIVE", "RESET_PASS", "RESET_EMAIL"}
)

func GetVerifyCode(c *gin.Context) {
	uid := c.MustGet("uid").(uint64)
	scope := c.Query("scope")
	if !isSupport(scope) {
		c.JSON(http.StatusBadRequest, model.NewResp().WithErrMsg(myerr.Auth, myerr.MsgIllegalVerifyScope))
		return
	}
	code := util.RandomNum(8)
	request_key := fmt.Sprintf("verify_code:%s:%d", scope, uid)
	err := db.GetRedis().Set(request_key, code, time.Minute*10).Err()
	if err != nil {
		log.Error().Err(err).Msg("generate verify code failed")
		c.JSON(http.StatusInternalServerError, myerr.ErrInternalServer)
		return
	}
	c.JSON(http.StatusOK, model.NewResp().AddResult("verify_code", code).AddResult("scope", scope))
}
func isSupport(scope string) bool {
	s := strings.ToUpper(scope)
	for _, v := range SUPPORT_VERIFY {
		if s == v {
			return true
		}
	}
	return false
}
