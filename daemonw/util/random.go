package util

import (
	"math/rand"
	"time"
)

const (
	Numbers = "0123456789"
)

func RandomBytes(len int) []byte {
	b := make([]byte, len)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len; i++ {
		b[i] = byte(rand.Intn(256));
	}
	return b
}

func RandomNum(len int) string {
	b := make([]byte, len)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < len; i++ {
		b[i] = Numbers[rand.Intn(256)%10];
	}
	return string(b)
}
