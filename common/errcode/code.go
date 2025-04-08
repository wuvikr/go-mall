package errcode

import "net/http"

var codes = map[int]string{}

var (
	Success            = newAppError(0, "success")
	ErrServer          = newAppError(10000000, "服务器内部错误")
	ErrParams          = newAppError(10000001, "参数错误, 请检查")
	ErrNotFound        = newAppError(10000002, "资源未找到")
	ErrPanic           = newAppError(10000003, "(*^__^*)系统开小差了,请稍后重试") // 无预期的panic错误
	ErrToken           = newAppError(10000004, "Token无效")
	ErrForbidden       = newAppError(10000005, "未授权") // 访问一些未授权的资源时的错误
	ErrTooManyRequests = newAppError(10000006, "请求过多")
)

func (e *AppError) HttpStatusCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	case ErrServer.Code():
		return http.StatusInternalServerError
	case ErrParams.Code():
		return http.StatusBadRequest
	case ErrNotFound.Code():
		return http.StatusNotFound
	case ErrTooManyRequests.Code():
		return http.StatusTooManyRequests
	case ErrToken.Code():
		return http.StatusUnauthorized
	case ErrForbidden.Code():
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
