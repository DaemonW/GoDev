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
	"github.com/asaskevich/govalidator"
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
	var registerUser struct {
		Username string `json:"username" form:"username" valid:"alphanum,length(8|16)"`
		Password string `json:"password" form:"password" valid:"printableascii,length(8|16)"`
	}
	if err = c.ShouldBindWith(&registerUser, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, model.NewResp().SetErrMsg(err.Error()))
		return
	}
	_, err = govalidator.ValidateStruct(registerUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewResp().SetErrMsg("bad format"))
		return
	}
	user := model.NewUser(registerUser.Username, registerUser.Password)
	if err = dao.UserDao.CreateUser(user); err != nil {
		resp := model.NewResp().SetError(model.ErrCreateUser)
		c.JSON(http.StatusBadRequest, resp)
	} else {
		resp := model.NewResp().AddResult("msg", "create user success")
		c.JSON(http.StatusOK, resp)
	}
}

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
		c.JSON(http.StatusBadRequest, model.NewResp().SetError(model.ErrUserNotExist))
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
			c.JSON(http.StatusInternalServerError, model.NewResp().SetError(model.ErrInternalServer))
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
		c.JSON(http.StatusUnauthorized, model.NewResp().SetError(model.ErrInvalidAuth))
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
