package middleware

import (
	"net/http"

	"github.com/bluele/gcache"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"daemonw/dao"
	"daemonw/model"
	"daemonw/db"
	"strconv"
)

var cache = gcache.New(1024).LRU().Build()

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusUnauthorized, model.NewResp().SetError(model.ErrInternalServer))
				c.Abort()
			}
		}()
		tokenStr := c.Request.Header.Get("auth")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, model.NewResp().SetErrMsg("invalid token"))
			c.Abort()
			return
		}
		token, err := verifyToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.NewResp().SetErrMsg(err.Error()))
			c.Abort()
			return
		}
		params := token.Claims.(jwt.MapClaims)
		c.Set("uid", params["uid"])
		c.Set("user", params["user"])
		c.Next()
	}
}

func verifyToken(tokenStr string) (*jwt.Token, error) {
	var pass string
	//check token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		uid := token.Claims.(jwt.MapClaims)["uid"].(uint64)
		//get cached pass
		result := db.GetRedis().Get("token_secret:" + strconv.FormatUint(uid, 10))
		//if cached, verified
		if result.Err() == nil {
			pass = result.String()
		}
		if pass != "" {
			return pass, nil
		}
		user, err := dao.NewUserDao().Get(uid)
		if err != nil {
			//internal error
			panic(err)
		}
		return user.Password, nil
	})
	return token, err
}
