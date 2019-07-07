package model

type UserInfo struct {
	UserId   uint64 `json:"userId" db:"user_id"`
	Nickname string `json:"nickname" db:"nickname"`
	Sex      uint8  `json:"sex" db:"sex"`
	Age      uint8  `json:"age" db:"age"`
	Email    string `json:"email" db:"email"`
	Phone    string `json:"phone" db:"phone"`
	Address  string `json:"address" db:"address"`
	Ip       string `json:"ip" db:"ip"`
	Meta     string `json:"meta" db:"meta"`
}
