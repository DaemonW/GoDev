package controller

import (
	"daemonw/errors"
	"net/http"
	"strconv"

	"daemonw/dao"
	"daemonw/log"
	"daemonw/model"
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
		c.JSON(http.StatusInternalServerError, errors.ErrInternalServer)
		return
	}
	if user == nil {
		c.JSON(http.StatusBadRequest, errors.ErrUserNotExist)
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
		c.JSON(http.StatusBadRequest, errors.ErrCreateUser)
	} else {
		resp := model.NewResp().AddResult("msg", "create user success")
		c.JSON(http.StatusOK, resp)
	}
}
