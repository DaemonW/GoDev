package controller

import (
	"daemonw/dao"
	"daemonw/entity"
	"daemonw/xerr"
	"daemonw/xlog"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
)

func CreateFile(c *gin.Context) {
	id,_ := strconv.ParseUint(c.Param("id"),10,64)
	file := entity.File{}
	if err := c.ShouldBindWith(&file, binding.JSON); err != nil {
		xlog.Error().Msgf("create file failed, reason: %s", err.Error())
		c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeCreateFile, xerr.MsgBadParam))
		return
	}
	file.Owner = id
	if file.ParentId < 0 {
		file.ParentId = 0;
	}
	fileDao:=dao.NewFileDao()
	err := fileDao.CreateFile(file)
	if err != nil {
		xlog.Error().Msgf("create file failed, reason: %s", err.Error())
		c.JSON(http.StatusInternalServerError, entity.NewRespErr(xerr.CodeCreateFile, xerr.MsgBadParam))
		return
	}
}
