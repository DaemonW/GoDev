package crypto

import (
	"math/rand"
	"time"
)

var (
	AlphaSeed = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLPNOPQRSTUVWXYZ!@#$%^&*()-_=+")
)

func RandomBytes(size int) []byte {
	b := make([]byte, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		b[i] = byte(rand.Uint32() & 0x0f)
	}
	return b
}

func RandomStr(size int) []byte {
	b := make([]byte, size)
	len := len(AlphaSeed)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		b[i] = AlphaSeed[byte(rand.Int()&0x7fffffff%len)]
	}
	return b
}
