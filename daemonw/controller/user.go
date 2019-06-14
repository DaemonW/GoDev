package controller

import (
	myerr "daemonw/errors"
	"net/http"
	"strconv"

	"crypto/tls"
	"daemonw/conf"
	"daemonw/dao"
	. "daemonw/model"
	"daemonw/util"
	"daemonw/xlog"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/gomail.v2"
	"time"
)

func GetUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	user, err := dao.UserDao.Get(id)
	if err != nil {
		xlog.Error().Err(err).Msg("db error")
		c.JSON(http.StatusInternalServerError, myerr.ErrInternalServer)
		return
	}
	if user == nil {
		c.JSON(http.StatusBadRequest, NewRespErr(myerr.QueryUser, myerr.MsgUserNotExist))
		return
	}
	c.JSON(http.StatusOK, NewResp().AddResult("user", user))
}

func GetUsers(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusOK, NewResp().AddResult("users", nil))
		return
	}
	users, err := dao.UserDao.GetLikeName(name)
	util.PanicIfErr(err)
	c.JSON(http.StatusOK, NewResp().AddResult("users", users))
}

func CreateUser(c *gin.Context) {
	var err error
	var registerUser struct {
		Username string `json:"username" valid:"email,length(4|32)"`
		Password string `json:"password" valid:"printableascii,length(8|16)"`
	}
	if err = c.ShouldBindWith(&registerUser, binding.JSON); err != nil {
		c.JSON(http.StatusNotAcceptable, NewRespErr(myerr.CreateUser, myerr.MsgBadParam))
		return
	}
	_, err = govalidator.ValidateStruct(registerUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewRespErr(myerr.CreateUser, myerr.MsgBadParam))
		return
	}
	user := NewUser(registerUser.Username, registerUser.Password)
	qUser, err := dao.UserDao.GetByName(registerUser.Username)
	if err != nil {
		xlog.Error().Msgf(err.Error())
		c.JSON(http.StatusInternalServerError, myerr.ErrInternalServer)
		return
	}
	if qUser != nil {
		c.JSON(http.StatusBadRequest, NewResp().WithErrMsg(myerr.CreateUser, myerr.MsgUserExist))
		return
	}
	if err = dao.UserDao.CreateUser(user); err != nil {
		xlog.Error().Msgf(err.Error())
		c.JSON(http.StatusBadRequest, NewRespErr(myerr.CreateUser, myerr.MsgCreateUserFail))
	} else {
		go sendMail(user)
		resp := NewResp().AddResult("msg", "create user success")
		c.JSON(http.StatusOK, resp)
	}
}

func ActiveUser(c *gin.Context) {
	id := c.Param("id")
	code := c.Query("code")
	if id == "" || code == "" {
		c.JSON(http.StatusBadRequest, NewResp().WithErrMsg(myerr.Login, myerr.MsgActiveUserFail))
		return
	}
	request_code := dao.Redis().Get("verify_code:active" + id).String()
	if request_code != code {
		c.JSON(http.StatusBadRequest, NewResp().WithErrMsg(myerr.Login, myerr.MsgActiveUserFail))
		return
	}
	uid, _ := strconv.ParseInt(id, 10, 64)
	err := dao.UserDao.ActiveUser(uid)
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, NewResp().AddResult("msg", "user is active"))
}

func SendActiveMail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {

	}
}

func sendMail(user *User) error {
	serverConf := conf.Config.SmtpServer
	d := gomail.NewDialer(serverConf.Host, serverConf.Port, serverConf.Account, serverConf.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	m := gomail.NewMessage()
	m.SetHeader("From", serverConf.Account)
	m.SetHeader("To", user.Username)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", genActiveUrl(user))
	m.Attach("/home/daemonw/Pictures/1531905129160.jpg")
	return d.DialAndSend(m)
}

func genActiveUrl(user *User) string {
	key := fmt.Sprintf("verify_code:active:%d", user.ID)
	code := util.RandomNum(16)
	err := dao.Redis().SetXX(key, code, time.Minute*10).Err()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("http://localhost:8080/api/user/%d?verify_code=%s", user.ID, code)
}
