package controller

import (
	myerr "daemonw/errors"
	"net/http"
	"strconv"

	"daemonw/dao"
	"daemonw/log"
	. "daemonw/model"
	"daemonw/util"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetUser(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("user_id"), 10, 64)
	user, err := dao.UserDao.Get(id)
	if err != nil {
		log.Error().Err(err).Msg("db error")
		c.JSON(http.StatusInternalServerError, myerr.ErrInternalServer)
		return
	}
	if user == nil {
		c.JSON(http.StatusBadRequest, NewRespErr(myerr.QueryUser,myerr.MsgUserNotExist))
		return
	}
	c.JSON(http.StatusOK, NewResp().AddResult("user", user))
}

func GetAllUsers(c *gin.Context) {
	users, err := dao.UserDao.GetAll()
	util.PanicIfErr(err)
	c.JSON(http.StatusOK, NewResp().AddResult("users", users))
}

func CreateUser(c *gin.Context) {
	var err error
	var registerUser struct {
		Username string `json:"username" form:"username" valid:"email,length(4|32)"`
		Password string `json:"password" form:"password" valid:"printableascii,length(8|16)"`
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
	if err = dao.UserDao.CreateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, NewRespErr(myerr.CreateUser,myerr.MsgCreateUserFail))
	} else {
		resp := NewResp().AddResult("msg", "create user success")
		c.JSON(http.StatusOK, resp)
	}
}
