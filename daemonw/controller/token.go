package controller

import (
	"crypto/md5"
	"daemonw/dao"
	. "daemonw/entity"
	"daemonw/util"
	"daemonw/xerr"
	"daemonw/xlog"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
	"time"
)

func Login(c *gin.Context) {
	var loginUser struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}
	if err := c.ShouldBindWith(&loginUser, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, NewRespErr(xerr.CodeLogin, err.Error()))
		return
	}

	u, err := dao.UserDao.GetByName(loginUser.Username)
	if err!=nil{
		c.JSON(http.StatusInternalServerError, NewRespErr(xerr.CodeInternal, xerr.MsgInternal))
		return
	}
	if u == nil {
		c.JSON(http.StatusBadRequest, NewRespErr(xerr.CodeLogin, xerr.MsgUserNotExist))
		return
	}

	if u.Status == StatusUserInactive {
		c.JSON(http.StatusBadRequest, NewRespErr(xerr.CodeLogin, xerr.MsgUserNotActive))
		return
	}

	if u.Status == StatusUserFreeze {
		c.JSON(http.StatusBadRequest, NewRespErr(xerr.CodeLogin, xerr.MsgFreezeUser))
		return
	}

	b := append([]byte(loginUser.Password), u.Salt...)
	encPass := fmt.Sprintf("%x", md5.Sum(b))
	if encPass == u.Password {
		ip := util.GetRequestIP(c.Request, false)
		xlog.Info().Msgf("%s request for login, ip = %s", loginUser.Username, ip)
		token, err := genJwtToken(u, ip)
		if err != nil {
			xlog.Error().Err(err).Msg("generate token failed")
			c.JSON(http.StatusInternalServerError, NewRespErr(xerr.CodeInternal, xerr.MsgInternal))
			return
		}
		c.Writer.Header().Set("auth", token)
		dao.Redis().Set("token_secret:"+strconv.FormatUint(u.ID, 10), u.Password, time.Minute*10)
		c.JSON(http.StatusOK,
			NewResp().
				AddResult("msg", "login success").
				AddResult("user", u).
				AddResult("token", token))
	} else {
		c.JSON(http.StatusUnauthorized, NewRespErr(xerr.CodeLogin, xerr.MsgIncorrectAuth))
	}
}

func genJwtToken(user *User, ip string) (string, error) {
	claims := Claims{
		Ip: ip,
		StandardClaims: jwt.StandardClaims{
			Id:        strconv.FormatUint(user.ID, 10),
			Issuer:    "server",
			ExpiresAt: time.Now().Add(time.Hour * 24 * 14).Unix(),
			Audience:  user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(user.Password))
}
