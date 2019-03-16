package middleware

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"daemonw/dao"
	"daemonw/model"
	"daemonw/db"
	myerr "daemonw/errors"
	"strconv"
	"daemonw/log"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error().Err(err.(error)).Msg("verify token")
				c.JSON(http.StatusInternalServerError, myerr.ErrInternalServer)
				c.Abort()
			}
		}()
		tokenStr := c.Request.Header.Get("auth")
		if tokenStr == "" {
			c.JSON(http.StatusUnauthorized, model.NewRespErr(myerr.Auth, "invalid token"))
			c.Abort()
			return
		}
		token, err := verifyToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, model.NewRespErr(myerr.Auth, err.Error()))
			c.Abort()
			return
		}
		params := token.Claims.(*model.Claims)
		uid, _ := strconv.ParseUint(params.Id, 10, 64)
		c.Set("uid", uid)
		c.Set("user", params.Audience)
		c.Next()
	}
}

func verifyToken(tokenStr string) (*jwt.Token, error) {
	var pass string
	var claims model.Claims
	//check token
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		claims, ok := token.Claims.(*model.Claims)
		if !ok {
			return nil, errors.New("invalid token format")
		}
		//get cached pass
		result, err := db.GetRedis().Get("token_secret:" + claims.Id).Result()
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
		db.GetRedis().Set("token_secret:"+claims.Id, user.Password, 0)
		return []byte(user.Password), nil
	})
	return token, err
}
