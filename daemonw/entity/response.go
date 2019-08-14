package entity

import "daemonw/xerr"

type Content map[string]interface{}

type Response struct {
	Result Content      `json:"result,omitempty"`
	Err    *xerr.Err `json:"err,omitempty"`
}

func NewResp() *Response {
	resp := &Response{}
	return resp
}

func NewRespErr(errCode int, errMsg string) *Response {
	resp := &Response{}
	resp.Err = &xerr.Err{errCode, errMsg}
	return resp
}

func (r *Response) AddResult(key string, val interface{}) *Response {
	if r.Result == nil {
		r.Result = make(map[string]interface{})
	}
	r.Result[key] = val
	return r
}

func (r *Response) WithErrMsg(errCode int, errMsg string) *Response {
	if r.Err == nil {
		r.Err = &xerr.Err{}
	}
	r.Err.Code = errCode;
	r.Err.Msg = errMsg
	return r
}

func (r *Response) SetError(err *xerr.Err) *Response {
	r.Err = err
	return r
}
