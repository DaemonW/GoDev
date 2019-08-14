package entity

import (
	"crypto/md5"
	"daemonw/util"
	"time"
)

const (
	UserStatusInactive = iota
	UserStatusNormal
	UserStatusFreeze

	UserRoleNormal = iota
	UserRoleAdmin
)

type User struct {
	Id       uint64    `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
	Password string    `json:"-" db:"password"`
	Salt     []byte    `json:"-" db:"salt"`
	Status   uint8     `json:"status" db:"status"`
	Role     uint8     `json:"role" db:"role"`
	CreateAt time.Time `json:"createAt" db:"create_at"`
	UpdateAt time.Time `json:"updateAt" db:"update_at"`
}

func NewUser(username, password string) *User {
	u := &User{Username: username, CreateAt: time.Now(), UpdateAt: time.Now()}
	u.Salt = util.RandomBytes(8)
	b := append([]byte(password), u.Salt...)
	hash := md5.Sum(b)
	encPass := util.Bytes2HexStr(hash[:])
	u.Password = encPass
	return u
}

func (u *User) SetPassword(password string, salt []byte) {
	u.Salt = salt
	if salt == nil || len(salt) != 8 {
		u.Salt = util.RandomBytes(8)
	}
	b := append([]byte(password), u.Salt...)
	hash := md5.Sum(b)
	encPass := util.Bytes2HexStr(hash[:])
	u.Password = encPass
}
