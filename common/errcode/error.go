package errcode

import (
	"encoding/json"
	"fmt"
	"path"
	"runtime"
)

type AppError struct {
	code     int
	msg      string
	cause    error
	occurred string
}

func newAppError(code int, msg string) *AppError {
	if _, ok := codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 已存在", code))
	}
	return &AppError{
		code: code,
		msg:  msg,
	}
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}

	formattedErr := struct {
		Code     int    `json:"code"`
		Msg      string `json:"msg"`
		Cause    string `json:"cause"`
		Occurred string `json:"occurred"`
	}{
		Code:     e.Code(),
		Msg:      e.Msg(),
		Occurred: e.occurred,
	}

	if e.cause != nil {
		formattedErr.Cause = e.cause.Error()
	}
	errBytes, _ := json.MarshalIndent(formattedErr, "", "  ")

	return string(errBytes)
}

func (e *AppError) Code() int {
	return e.code
}

func (e *AppError) Msg() string {
	return e.msg
}

func (e *AppError) String() string {
	return e.Error()
}

func Wrap(msg string, err error) *AppError {
	if err == nil {
		return nil
	}

	return &AppError{
		code:     -1,
		msg:      msg,
		cause:    err,
		occurred: getAppErrOccurredInfo(),
	}

}

func getAppErrOccurredInfo() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	file = path.Base(file)
	funcName := runtime.FuncForPC(pc).Name()
	triggerInfo := fmt.Sprintf("func:%s, file:%s, line:%d", funcName, file, line)
	return triggerInfo
}

// WithCause 在逻辑执行中出现错误, 比如dao层返回的数据库查询错误
// 可以在领域层返回预定义的错误前附加上导致错误的基础错误。
// 如果业务模块预定义的错误码比较详细, 可以使用这个方法, 反之错误码定义的比较笼统建议使用Wrap方法包装底层错误生成项目自定义Error
// 并将其记录到日志后再使用预定义错误码返回接口响应
func (e *AppError) WithCause(err error) *AppError {
	e.cause = err
	e.occurred = getAppErrOccurredInfo()
	return e
}
