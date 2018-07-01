package controller

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strconv"

	"daemonw/dao"
	"daemonw/model"
	"daemonw/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/dgrijalva/jwt-go"
	"daemonw/log"
	"daemonw/db"
	"time"
)

func GetUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	user, err := dao.UserDao.Get(id)
	if err != nil {
		log.Error().Err(err).Msg("db error")
		c.JSON(http.StatusInternalServerError, model.NewResp().SetError(model.ErrInternalServer))
		return
	}
	if user == nil {
		c.JSON(http.StatusBadRequest, model.NewResp().SetError(model.ErrUserNotExist))
		return
	}
	c.JSON(http.StatusOK, model.NewResp().AddResult("user", user))
}

func GetAllUsers(c *gin.Context) {
	users, err := dao.UserDao.GetAll()
	util.PanicIfErr(err)
	c.JSON(http.StatusOK, model.NewResp().AddResult("users", users))
}

func CreateUser(c *gin.Context) {
	var err error
	var user model.User
	if err = c.ShouldBindWith(&user, binding.FormPost); err != nil {
		c.JSON(http.StatusBadRequest, model.NewResp().SetErrMsg(err.Error()))
		return
	}
	user.SetPassword(user.Password, nil)
	if err = dao.UserDao.CreateUser(&user); err != nil {
		resp := model.NewResp().SetError(model.ErrCreateUser)
		c.JSON(http.StatusBadRequest, resp)
	} else {
		resp := model.NewResp().AddResult("msg", "create user success")
		c.JSON(http.StatusOK, resp)
	}
}

func Login(c *gin.Context) {
	var user struct {
		Username string `json:"username" form:"username"`
		Password string `json:"password" form:"password"`
	}
	if err := c.ShouldBindWith(&user, binding.FormPost); err != nil {
		c.JSON(http.StatusBadRequest, model.NewResp().SetErrMsg(err.Error()))
		return
	}

	u, err := dao.UserDao.GetByName(user.Username)
	util.PanicIfErr(err)
	if u == nil {
		c.JSON(http.StatusBadRequest, model.NewResp().SetError(model.ErrUserNotExist))
		return
	}

	b := append([]byte(user.Password), u.Salt...)
	encPass := fmt.Sprintf("%x", md5.Sum(b))
	if encPass == u.Password {
		ip := util.GetRequestIP(c.Request, false)
		u.LoginIp = ip
		token, err := genJwtToken(u, ip)
		if err != nil {
			log.Error().Err(err).Msg("generate token failed")
			c.JSON(http.StatusInternalServerError, model.NewResp().SetError(model.ErrInternalServer))
			return
		}
		c.Writer.Header().Set("auth", token)
		db.GetRedis().Set("token_secret:"+strconv.FormatUint(u.ID, 10), u.Password, time.Minute*10)
		c.JSON(http.StatusOK,
			model.NewResp().
				AddResult("msg", "login success, ip address = "+ip).
				AddResult("user", u))
	} else {
		c.JSON(http.StatusUnauthorized, model.NewResp().SetError(model.ErrInvalidAuth))
	}
}

func genJwtToken(user *model.User, ip string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":    strconv.FormatUint(user.ID, 10),
		"user":   user.Username,
		"create": user.CreateAt,
		"ip":     ip,
	})
	return token.SignedString([]byte(user.Password))
}
