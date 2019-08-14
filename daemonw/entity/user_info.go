package entity

type UserInfo struct {
	Id   uint64 `json:"id" db:"id"`
	Nickname string `json:"nickname" db:"nickname"`
	Sex      uint8  `json:"sex" db:"sex"`
	Age      uint8  `json:"age" db:"age"`
	Email    string `json:"email" db:"email"`
	Phone    string `json:"phone" db:"phone"`
	Address  string `json:"address" db:"address"`
	Ip       string `json:"ip" db:"ip"`
	Extra     string `json:"extra" db:"extra"`
}
