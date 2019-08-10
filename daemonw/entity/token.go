package entity

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Ip string `json:"ip"`
	jwt.StandardClaims
}
