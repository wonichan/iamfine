package snerror

import "fmt"

type SnError struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Err  error
}

func NewSnError(code int32, msg string) SnError {
	return SnError{
		Code: code,
		Msg:  msg,
	}
}

func NewSnErrorWithError(code int32, err error) SnError {
	return SnError{
		Code: code,
		Err:  err,
	}
}

func (s SnError) Error() string {
	if s.Err != nil {
		return fmt.Sprintf("%d:%s", s.Code, s.Err.Error())
	}
	return fmt.Sprintf("%d:%s", s.Code, s.Msg)
}
