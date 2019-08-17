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
)

const (
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
	u.Password, u.Salt = GenPassword(password, nil)
	u.CreateAt = time.Now()
	u.UpdateAt = time.Now()
	return u
}

func GenPassword(password string, salt []byte) (string, []byte) {
	var s []byte
	if salt == nil || len(salt) != 8 {
		s = util.RandomBytes(8)
	}
	b := append([]byte(password), s...)
	hash := md5.Sum(b)
	p := util.Bytes2HexStr(hash[:])
	return p, s
}
