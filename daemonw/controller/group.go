package controller

import (
	"daemonw/dao"
	"daemonw/entity"
	"daemonw/xerr"
	"daemonw/xlog"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	"strconv"
	"strings"
)

func CreateGroup(c *gin.Context) {
	t := c.Request.Header.Get("Content-Type")
	if !strings.HasPrefix(t, gin.MIMEJSON) {
		fmt.Println(t)
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeCreateGroup, "unacceptable content type"))
		return
	}
	group := entity.AppGroup{}
	err := c.BindJSON(&group)
	if err != nil {
		msg := fmt.Sprintf("parse params failed: %s", err.Error())
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeCreateGroup, msg))
		return
	}
	err = dao.NewGroupDao().CreateGroup(&group)
	exist := false
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			code := pqErr.Code
			if code == "23505" {
				exist = true
			}
		}
		xlog.Error().Msgf("insert group into db failed: %s", err.Error())
		msg := "insert group failed"
		if exist {
			msg = "insert group failed: group exists"
		}
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeCreateGroup, msg))
		return
	}
	c.JSON(http.StatusOK, entity.NewResp().AddResult("group", group))
}

func DeleteGroup(c *gin.Context) {
	idStr := c.Param("id")
	_id, _ := strconv.ParseUint(idStr, 10, 64)
	if _id <= 0 {
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeDeleteGroup, "illegal id param"))
		return
	}
	_, err := dao.NewGroupDao().Delete(`DELETE FROM groups WHERE id = ?`, _id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.NewRespErr(xerr.CodeDeleteGroup, err.Error()))
	}
	c.JSON(http.StatusOK, entity.NewResp().AddResult("msg", "delete group success"))
}

func GetGroups(c *gin.Context) {
	groups, err := dao.NewGroupDao().GetAllGroups()
	if err != nil {
		xlog.Error().Msgf("query groups failed: %s", err.Error())
		c.JSON(http.StatusInternalServerError, entity.NewRespErr(xerr.CodeInternal, xerr.MsgInternal))
		return
	}
	c.JSON(http.StatusOK, entity.NewResp().AddResult("groups", groups))
}

func QueryGroupApps(c *gin.Context) {
	apps := make([]entity.App, 0)
	dao.NewGroupDao().Select(&apps)
}
