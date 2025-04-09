package app

import (
	"go-mall/common/errcode"
	"go-mall/common/logger"

	"github.com/gin-gonic/gin"
)

type response struct {
	ctx        *gin.Context
	Code       int         `json:"code"`
	Msg        any         `json:"msg"`
	RequestId  any         `json:"request_id"`
	Data       any         `json:"data,omitempty"`
	Pagination *pagination `json:"pagination,omitempty"`
}

func NewResponse(ctx *gin.Context) *response {
	return &response{
		ctx: ctx,
	}
}

func (r *response) SetPagination(p *pagination) *response {
	r.Pagination = p
	return r
}

func (r *response) Success(data any) {
	r.Code = errcode.Success.Code()
	r.Msg = errcode.Success.Msg()
	requestId := ""
	if _, exists := r.ctx.Get("traceid"); exists {
		requestId = r.ctx.GetString("traceid")
	}
	r.RequestId = requestId
	r.Data = data

	r.ctx.JSON(errcode.Success.HttpStatusCode(), r)
}

func (r *response) SuccessOk() {
	r.Success("")
}

func (r *response) Error(err *errcode.AppError) {
	r.Code = err.Code()
	r.Msg = err.Msg()
	requestId := ""
	if _, exists := r.ctx.Get("traceid"); exists {
		requestId = r.ctx.GetString("traceid")
	}
	r.RequestId = requestId

	// 兜底错误log
	logger.New(r.ctx).Error("api_response_error", "err", err)
	r.ctx.JSON(err.HttpStatusCode(), r)
}
