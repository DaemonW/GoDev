package xerr

import "fmt"

type ErrMsg struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg"`
}

func (err *ErrMsg) String() string {
	return fmt.Sprintf("errno: %d, %s", err.Code, err.Msg)
}

func (err *ErrMsg) Error() string {
	return fmt.Sprintf("errno: %d, %s", err.Code, err.Msg)
}
