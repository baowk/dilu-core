package errs

import (
	"errors"
	"fmt"
)

type ErrType string

const (
	DB_ERR ErrType = "DB_ERR"
)

func (e ErrType) String() string {
	return string(e)
}

type BizError struct {
	//error
	reqId string
	code  int
	//msg       string
	data      map[string]any
	causes    []error
	retryable bool
}
type IError interface {
	error
	Code() int
	//Msg() string
	ReqId() string
	Causes() []error
	Data() map[string]interface{}
	Retryable() bool
	SetRetryable(r bool)
}

func (e *BizError) Error() string {
	return fmt.Sprintf("[reqId]:%s [Code]:%d [data]:%v [errors] %s", e.reqId, e.code, e.data, e.causes)
}

func (e *BizError) Causes() []error {
	return e.causes
}

func (e *BizError) Data() map[string]any {
	return e.data
}

func (e *BizError) Code() int {
	return e.code
}

func (e *BizError) ReqId() string {
	return e.reqId
}

// func (e *BizError) Msg() string {
// 	return e.msg
// }

func (e *BizError) Retryable() bool {
	return e.retryable
}

func (e *BizError) SetRetryable(r bool) {
	e.retryable = r
}

func (e *BizError) GetDataVal(key string) any {
	return e.data[key]
}

func Err(code int, reqId string, cause error) IError {
	return &BizError{
		reqId:  reqId,
		code:   code,
		causes: []error{cause},
	}
}

func ErrWithData(code int, reqId string, cause error, m map[string]any) IError {
	return &BizError{
		reqId:  reqId,
		code:   code,
		causes: []error{cause},
		data:   m,
	}
}

func AsBizError(err error) *BizError {
	var bizError = new(BizError)
	if errors.As(err, &bizError) {
		return bizError
	}
	return nil
}
