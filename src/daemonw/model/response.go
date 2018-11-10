package model

import "daemonw/errors"

type Content map[string]interface{}

type Response struct {
	Result Content        `json:"result,omitempty"`
	Err    *errors.ErrMsg `json:"err,omitempty"`
}

func NewResp() *Response {
	resp := &Response{}
	return resp
}

func (r *Response) AddResult(key string, val interface{}) *Response {
	if r.Result == nil {
		r.Result = make(map[string]interface{})
	}
	r.Result[key] = val
	return r
}

func (r *Response) SetErrCode(errCode int) *Response {
	if r.Err == nil {
		r.Err = &errors.ErrMsg{}
	}
	r.Err.Code = errCode
	return r
}

func (r *Response) SetErrMsg(errMsg string) *Response {
	if r.Err == nil {
		r.Err = &errors.ErrMsg{}
	}
	r.Err.Msg = errMsg
	return r
}

func (r *Response) SetError(err *errors.ErrMsg) *Response {
	r.Err = err
	return r
}
