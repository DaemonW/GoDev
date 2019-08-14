package xerr

import "fmt"

type Err struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg"`
}

func (err *Err) String() string {
	return fmt.Sprintf("errno: %d, %s", err.Code, err.Msg)
}

func (err *Err) Error() string {
	return fmt.Sprintf("errno: %d, %s", err.Code, err.Msg)
}

func (err *Err) IsInternalErr() bool {
	return err.Code < CodeBiz
}
