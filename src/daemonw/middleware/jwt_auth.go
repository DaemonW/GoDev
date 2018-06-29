package middleware

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"daemonw/dao"
	"daemonw/model"
	"daemonw/db"
	"strconv"
	"errors"
	"daemonw/log"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, model.NewResp().SetError(model.ErrInternalServer))
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
		uidStr := params["uid"].(string)
		uid, _ := strconv.ParseUint(uidStr, 10, 64)
		c.Set("uid", uid)
		c.Set("user", params["user"])
		c.Next()
	}
}

func verifyToken(tokenStr string) (*jwt.Token, error) {
	var pass string
	//check token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		uidStr := token.Claims.(jwt.MapClaims)["uid"]
		//get cached pass
		result, err := db.GetRedis().Get("token_secret:" + uidStr.(string)).Result()
		//if cached, verified
		if err == nil {
			pass = result
		}
		if pass != "" {
			log.Info().Msgf("password = %s", pass)
			return []byte(pass), nil
		}
		uid, _ := strconv.ParseUint(uidStr.(string), 10, 64)
		user, err := dao.UserDao.Get(uid)
		if err != nil {
			//internal error
			panic(err)
		}
		if user == nil {
			return nil, errors.New("invalid token")
		}
		return []byte(user.Password), nil
	})
	return token, err
}
