package controller

import (
	"github.com/gin-gonic/gin"
	"daemonw/middleware"
	"daemonw/db"
	"strconv"
	"time"
	"net/http"
	"daemonw/model"
	"daemonw/conf"
	"daemonw/log"
	"errors"
)

func LimitUserAccessCount(c *gin.Context) {
	uid := c.MustGet("uid").(uint64)
	key := conf.AccessCountLimitKey + strconv.FormatUint(uid, 10)
	beginStr := c.Request.FormValue("begin")
	endStr := c.Request.FormValue("begin")
	limitStr := c.Request.FormValue("limit")
	limiter := middleware.NewCounter(db.GetRedis())
	if limitStr == "" {
		limiter.DeleleLimit(key)
	}
	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewResp().SetErrMsg("invalid params"))
		return
	}
	begin, end, err := parseTime(beginStr, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.NewResp().SetErrMsg(err.Error()))
		return
	}
	limiter.SetLimit(key, 10, begin, end)
	log.Debug().Msgf("set count limit, limit count = %d, start = %s, end = %s",
		limit,
		begin.Format("2006-01-02 15:04:05"),
		end.Format("2006-01-02 15:04:05"))
	c.JSON(http.StatusOK, model.NewResp().AddResult("msg", "set time success"))
}

func parseTime(begin, end string) (beginTime time.Time, endTime time.Time, err error) {
	var unixBegin, unixEnd int64

	if end == "" {
		endTime = time.Now().Add(time.Second * 60)
	} else {
		unixEnd, err = strconv.ParseInt(end, 10, 64)
		if err != nil {
			return
		}
		endTime = time.Unix(unixEnd, 0)
	}

	if begin == "" {
		beginTime = time.Now()
	} else {
		unixBegin, err = strconv.ParseInt(end, 10, 64)
		if err != nil {
			return
		}
		beginTime = time.Unix(unixBegin, 0)
	}

	if endTime.Before(beginTime) {
		err = errors.New("invalid params")
	}
	return
}
