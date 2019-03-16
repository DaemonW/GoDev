package model

import (
	"daemonw/util"
	"crypto/md5"
	"fmt"
	"time"
)

const (
	STATUS_UNACTIVE = iota
	STATUS_NORMAL
	STATUS_FREEZE
)

type User struct {
	ID       uint64    `json:"id" db:"id"`
	Username string    `json:"username"`
	Password string    `json:"-"`
	Salt     []byte    `json:"-"`
	Status   uint8     `json:"status"`
	CreateAt time.Time `json:"createAt" db:"create_at"`
	UpdateAt time.Time `json:"updateAt" db:"update_at"`
}

func NewUser(username, password string) *User {
	u := &User{Username: username, CreateAt: time.Now(), UpdateAt: time.Now()}
	u.Salt = util.RandomBytes(8)
	b := append([]byte(password), u.Salt...)
	encPass := fmt.Sprintf("%x", md5.Sum(b))
	u.Password = encPass
	return u
}

func (u *User) SetPassword(password string, salt []byte) {
	u.Salt = salt
	if salt == nil || len(salt) != 8 {
		u.Salt = util.RandomBytes(8)
	}
	b := append([]byte(password), u.Salt...)
	encPass := fmt.Sprintf("%x", md5.Sum(b))
	u.Password = encPass
}
