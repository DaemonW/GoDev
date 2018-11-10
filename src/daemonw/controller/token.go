package controller

import (
	"crypto/md5"
	"daemonw/dao"
	"daemonw/db"
	"daemonw/errors"
	"daemonw/log"
	"daemonw/model"
	"daemonw/util"
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
	if err := c.ShouldBindWith(&loginUser, binding.FormPost); err != nil {
		c.JSON(http.StatusBadRequest, model.NewResp().SetErrMsg(err.Error()))
		return
	}

	u, err := dao.UserDao.GetByName(loginUser.Username)
	util.PanicIfErr(err)
	if u == nil {
		c.JSON(http.StatusBadRequest, errors.ErrUserNotExist)
		return
	}

	b := append([]byte(loginUser.Password), u.Salt...)
	encPass := fmt.Sprintf("%x", md5.Sum(b))
	if encPass == u.Password {
		ip := util.GetRequestIP(c.Request, false)
		u.LoginIp = ip
		token, err := genJwtToken(u, ip)
		if err != nil {
			log.Error().Err(err).Msg("generate token failed")
			c.JSON(http.StatusInternalServerError, errors.ErrInternalServer)
			return
		}
		c.Writer.Header().Set("auth", token)
		db.GetRedis().Set("token_secret:"+strconv.FormatUint(u.ID, 10), u.Password, time.Minute*10)
		c.JSON(http.StatusOK,
			model.NewResp().
				AddResult("msg", "login success, ip address = "+ip).
				AddResult("user", u).
				AddResult("token", token))
	} else {
		c.JSON(http.StatusUnauthorized, errors.ErrInvalidAuth)
	}
}

func genJwtToken(user *model.User, ip string) (string, error) {
	claims := model.Claims{
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
