package errors

import "fmt"

var (
	ErrInternalServer = &ErrMsg{Internal, MsgInternal,}
	ErrUserNotExist   = &ErrMsg{CreateUser, MsgUserNotExist}
	ErrInvalidAuth    = &ErrMsg{Login,MsgIncorrectAuth}
	ErrCreateUser     = &ErrMsg{CreateUser, MsgCreateUser}
)

type ErrMsg struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg"`
}

func (err *ErrMsg) String() string {
	return fmt.Sprintf("err = %s, code = %d", err.Msg, err.Code)
}
