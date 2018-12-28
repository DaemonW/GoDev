package model

type UserInfo struct {
	ID       uint64 `json:"id",db:"id"`
	UserId   uint64 `json:"userId",db:"user_id"`
	Nickname string `json:"nickname",db:"nickname"`
	Sex      uint8  `json:"sex",db:"sex"`
	Email    string `json:"email",db:"email"`
	Phone    string `json:"phone",db:"phone"`
	Address  string `json:"address",db:"address"`
	IpAddr   string `json:"ipAddr",db:"ip_addr"`
}
