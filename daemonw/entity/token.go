package entity

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	Ip   string `json:"ip"`
	Role uint8 `json:"role"`
	jwt.StandardClaims
}
