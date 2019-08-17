package middleware

import (
	"errors"
	"net/http"

	"daemonw/dao"
	"daemonw/entity"
	"daemonw/xerr"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"strconv"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.Request.Header.Get("auth")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, entity.NewRespErr(xerr.CodeAuth, "invalid token"))
			c.Abort()
			return
		}
		token, err := verifyToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, entity.NewRespErr(xerr.CodeAuth, err.Error()))
			c.Abort()
			return
		}
		params := token.Claims.(*entity.Claims)
		uid, err:= strconv.ParseUint(params.Id, 10, 64)
		if err!=nil{
			c.JSON(http.StatusBadRequest, entity.NewRespErr(xerr.CodeAuth, err.Error()))
			c.Abort()
			return
		}
		c.Set("uid", uid)
		c.Set("user", params.Audience)
		c.Next()
	}
}

func verifyToken(tokenStr string) (*jwt.Token, error) {
	var pass string
	claims:= entity.Claims{}
	//check token
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		//get cached pass
		result, err := dao.Redis().Get("token_secret:" + claims.Id).Result()
		//if cached, verified
		if err == nil {
			pass = result
		}
		if pass != "" {
			return []byte(pass), nil
		}
		uid, _ := strconv.ParseUint(claims.Id, 10, 64)
		user, err := dao.UserDao.Get(uid)
		if err != nil {
			//internal error
			panic(err)
		}
		if user == nil {
			return nil, errors.New("invalid token")
		}
		dao.Redis().Set("token_secret:"+claims.Id, user.Password, 0)
		return []byte(user.Password), nil
	})

	claims, ok := token.Claims.(entity.Claims)
	if !ok {
		return nil, errors.New("invalid token format")
	}
	return token, err
}
