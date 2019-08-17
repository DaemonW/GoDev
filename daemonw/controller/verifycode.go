package controller

import (
	"daemonw/dao"
	"daemonw/entity"
	"daemonw/util"
	"daemonw/xerr"
	"daemonw/xlog"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	SupportVerify = []string{"ACTIVE_USER", "RESET_PASS", "RESET_EMAIL"}
)

const (
	CodeExpire = time.Minute * 10;
)

func GetVerifyCode(c *gin.Context) {
	scope := c.Query("scope")
	if !isSupport(scope) {
		c.JSON(http.StatusBadRequest, entity.NewResp().WithErrMsg(xerr.CodeVerify, xerr.MsgIllegalVerifyScope))
		return
	}
	uid,_ := strconv.ParseUint(c.Query("uid"),10,64)
	if uid<=0{
		c.JSON(http.StatusBadRequest, entity.NewResp().WithErrMsg(xerr.CodeVerify, xerr.MsgUserNotExist))
		return
	}
	code := util.RandomNum(8)
	requestKey := fmt.Sprintf("verify_code:%s:%d", scope, uid)
	err := dao.Redis().Set(requestKey, code, CodeExpire).Err()
	if err != nil {
		xlog.Error().Err(err).Msg("generate verify code failed")
		c.JSON(http.StatusInternalServerError, entity.NewRespErr(xerr.CodeInternal, xerr.MsgInternal))
		return
	}
	c.JSON(http.StatusOK, entity.NewResp().AddResult("verify_code", code).AddResult("scope", scope))
}
func isSupport(scope string) bool {
	s := strings.ToUpper(scope)
	for _, v := range SupportVerify {
		if s == v {
			return true
		}
	}
	return false
}
