package controller

import (
	"daemonw/xerr"
	"net/http"
	"strconv"

	"crypto/tls"
	"daemonw/conf"
	"daemonw/dao"
	. "daemonw/entity"
	"daemonw/util"
	"daemonw/xlog"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gopkg.in/gomail.v2"
	"time"
)

func GetUsers(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusOK, NewResp().AddResult("users", nil))
		return
	}
	if name == "*" {
		role := c.MustGet("role").(uint8)
		if role == UserRoleNormal {
			c.JSON(http.StatusOK, NewResp().AddResult("users", nil))
		} else if role == UserRoleAdmin {
			users, err := dao.UserDao.GetAll()
			if err != nil {
				xlog.Panic().Msg(util.StackInfo())
				c.JSON(http.StatusInternalServerError, NewRespErr(xerr.CodeInternal, xerr.MsgInternal))
				return
			}
			c.JSON(http.StatusOK, NewResp().AddResult("users", users))
		}
		return
	}
	users, err := dao.UserDao.GetLikeName(name)
	if err != nil {
		xlog.Panic().Msg(util.StackInfo())
		c.JSON(http.StatusInternalServerError, NewRespErr(xerr.CodeInternal, xerr.MsgInternal))
		return
	}
	c.JSON(http.StatusOK, NewResp().AddResult("users", users))
}

func GetUser(c *gin.Context) {
	_id := c.Param("id")
	id, _ := strconv.ParseUint(_id, 10, 64);
	if id <= 0 {
		c.JSON(http.StatusNotFound, NewRespErr(xerr.CodeQueryUser, xerr.MsgUserNotExist))
		return
	}

	users, err := dao.UserDao.Get(id)
	if err != nil {
		xlog.Panic().Msg(util.StackInfo())
		c.JSON(http.StatusInternalServerError, NewRespErr(xerr.CodeInternal, xerr.MsgInternal))
		return
	}
	c.JSON(http.StatusOK, NewResp().AddResult("users", users))
}

func CreateUser(c *gin.Context) {
	var err error
	var registerUser struct {
		Username string `json:"username" valid:"email,length(4|32)"`
		Password string `json:"password" valid:"printableascii,length(8|16)"`
	}
	if err = c.ShouldBindWith(&registerUser, binding.JSON); err != nil {
		c.JSON(http.StatusNotAcceptable, NewRespErr(xerr.CodeCreateUser, xerr.MsgBadParam))
		return
	}
	_, err = govalidator.ValidateStruct(registerUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewRespErr(xerr.CodeCreateUser, xerr.MsgBadParam))
		return
	}
	user := NewUser(registerUser.Username, registerUser.Password)
	qUser, err := dao.UserDao.GetByName(registerUser.Username)
	if err != nil {
		xlog.Panic().Msg(util.StackInfo())
		c.JSON(http.StatusInternalServerError, NewRespErr(xerr.CodeInternal, xerr.MsgInternal))
		return
	}
	if qUser != nil {
		c.JSON(http.StatusBadRequest, NewResp().WithErrMsg(xerr.CodeCreateUser, xerr.MsgUserExist))
		return
	}
	if err = dao.UserDao.CreateUser(user); err != nil {
		xlog.Panic().Msg(util.StackInfo())
		c.JSON(http.StatusBadRequest, NewRespErr(xerr.CodeCreateUser, xerr.MsgCreateUserFail))
	} else {
		go sendMail(user)
		resp := NewResp().AddResult("msg", "create user success")
		c.JSON(http.StatusOK, resp)
	}
}

func UpdateUser(c *gin.Context) {
	role := c.MustGet("role").(uint8)
	_status := c.Query("status")
	if !util.IsEmpty(_status) {
		code := c.Query("code")
		status, _ := strconv.ParseUint(_status, 10, 8)
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		if id == 0 {
			c.JSON(http.StatusBadRequest, NewRespErr(xerr.CodeUpdateUser, xerr.MsgUserNotExist))
			return
		}

		if util.IsEmpty(code) {
			c.JSON(http.StatusBadRequest, NewRespErr(xerr.CodeUpdateUser, xerr.MsgIllegalRequestCode))
			return
		}
		u, err := dao.UserDao.Get(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, NewRespErr(xerr.CodeInternal, xerr.MsgInternal))
			return
		}
		if err1 := updateUserStatus(u, uint8(status), code, role == UserRoleAdmin); err1 != nil {
			if err1.IsInternalErr() {
				c.JSON(http.StatusInternalServerError, err1)
			} else {
				c.JSON(http.StatusBadRequest, err1)
			}
		}
		return
	}
}

func updateUserStatus(user *User, newStatus uint8, code string, isAdmin bool) *xerr.Err {
	if user==nil{
		return nil
	}

	if user.Status == newStatus {
		return nil
	}

	if user.Status == UserStatusNormal {
		if !isAdmin {
			return &xerr.Err{xerr.CodeUpdateUser, xerr.MsgPermissionDenied}
		}
	} else if user.Status == UserStatusFreeze {
		if !isAdmin {
			return &xerr.Err{xerr.CodeUpdateUser, xerr.MsgPermissionDenied}
		}

	} else if user.Status == UserStatusInactive {
		if !isAdmin {
			if newStatus == UserStatusNormal {
				requestCode := dao.Redis().Get("verify_code:active" + strconv.FormatUint(user.Id, 10)).String()
				if requestCode != code {
					return &xerr.Err{xerr.CodeUpdateUser, xerr.MsgIllegalRequestCode}
				}
			} else {
				return &xerr.Err{xerr.CodeUpdateUser, xerr.MsgPermissionDenied}
			}
		}
	}
	err := dao.UserDao.UpdateStatus(user.Id, newStatus)
	if err != nil {
		xlog.Panic().Msg(util.StackInfo())
		return &xerr.Err{xerr.CodeInternal, xerr.MsgInternal}
	}
	return nil
}

func DeleteUser(c *gin.Context) {
	role, _ := c.MustGet("role").(uint8)
	if role != UserRoleAdmin {
		c.JSON(http.StatusBadRequest, NewRespErr(xerr.CodeDelUser, xerr.MsgPermissionDenied))
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 0, 64)
	if id <= 0 {
		return
	}
	err := dao.UserDao.DeleteUser(id)
	if err != nil {
		xlog.Panic().Msg(util.StackInfo())
		c.JSON(http.StatusInternalServerError, NewRespErr(xerr.CodeInternal, xerr.MsgInternal))
		return
	}
}

func sendMail(user *User) {
	serverConf := conf.Config.SMTPServer
	d := gomail.NewDialer(serverConf.Host, serverConf.Port, serverConf.Account, serverConf.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	m := gomail.NewMessage()
	m.SetHeader("From", serverConf.Account)
	m.SetHeader("To", user.Username)
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", genActiveUrl(user))
	//m.Attach("/home/daemonw/test.jpg")
	err := d.DialAndSend(m)
	if err != nil {
		xlog.Error().Msgf("send email failed, user: %s, reason: %s", user.Username, err.Error())
	}
}

func genActiveUrl(user *User) string {
	key := fmt.Sprintf("verify_code:active:%d", user.Id)
	code := util.RandomNum(16)
	err := dao.Redis().SetXX(key, code, time.Minute*10).Err()
	if err != nil {
		xlog.Panic().Msg(util.StackInfo())
	}
	return fmt.Sprintf("http://%s:%d/api/user/%d?verify_code=%s", conf.Config.Domain, conf.Config.Port, user.Id, code)
}
