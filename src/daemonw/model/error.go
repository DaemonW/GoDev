package model

import "fmt"

var (
	ErrInternalServer = &ErrMsg{Msg: "internal server error"}
	ErrUserNotExist    = &ErrMsg{Msg: "user not exists"}
	ErrInvalidAuth    = &ErrMsg{Msg: "incorrect username or password"}
	ErrCreateUser    = &ErrMsg{Msg: "failed to create user"}
)

type ErrMsg struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg"`
}

func (err *ErrMsg) String() string {
	return fmt.Sprintf("err = %s, code = %d", err.Msg, err.Code)
}
