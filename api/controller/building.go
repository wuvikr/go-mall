package controller

import (
	"errors"
	"go-mall/api/request"
	"go-mall/common/app"
	"go-mall/common/errcode"
	"go-mall/common/logger"
	"go-mall/config"
	"go-mall/logic/appservice"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TestPing(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"msg": "pong",
	})
}

func TestConfigRead(c *gin.Context) {
	database := config.Database
	logger.ZapLoggerTest(c)
	log := config.App.Log
	c.JSON(http.StatusOK, gin.H{
		"type":        database.Type,
		"log.path":    log.FilePath,
		"log.maxsize": log.FileMaxSize,
		"log.maxage":  log.BackUpFileMaxAge,
	})
}

func TestLogger(c *gin.Context) {
	logger.New(c).Info("logger-test", "key1", "value1", "key2", 2)

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func TestAccessLog(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func TestPanicLog(c *gin.Context) {
	var a map[string]string
	a["k"] = "v"
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"data":   a,
	})
}

func TestAppError(c *gin.Context) {

	// 使用 Wrap 包装原因error 生成 项目error
	err := errors.New("a dao error")
	appErr := errcode.Wrap("包装错误", err)
	bAppErr := errcode.Wrap("再包装错误", appErr)
	logger.New(c).Error("记录错误", "err", bAppErr)

	// 预定义的ErrServer, 给其追加错误原因的error
	err = errors.New("a domain error")
	apiErr := errcode.ErrServer.WithCause(err)
	logger.New(c).Error("API执行中出现错误", "err", apiErr)

	c.JSON(apiErr.HttpStatusCode(), gin.H{
		"code": apiErr.Code(),
		"msg":  apiErr.Msg(),
	})
}

func TestResponseObj(c *gin.Context) {
	data := map[string]string{
		"key":  "value",
		"key2": "value2",
	}
	app.NewResponse(c).Success(data)
}

func TestResponseError(c *gin.Context) {
	baseErr := errors.New("a dao error")
	err := errcode.Wrap("包装错误", baseErr)
	app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
}

func TestResponseList(ctx *gin.Context) {
	pagination := app.NewPagination(ctx)

	data := []struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{
		{
			Name: "zhangsan",
			Age:  22,
		},
		{
			Name: "lisi",
			Age:  15,
		},
	}
	pagination.SetTotalRows(2)
	app.NewResponse(ctx).SetPagination(pagination).Success(data)

}

func TestGormLogger(ctx *gin.Context) {
	svc := appservice.NewDemoAppSvc(ctx)
	list, err := svc.GetDemoIdentities()
	if err != nil {
		app.NewResponse(ctx).Error(errcode.ErrServer.WithCause(err))
		return
	}
	app.NewResponse(ctx).Success(list)
}

func TestCreateDemoOrder(c *gin.Context) {
	request := new(request.DemoOrderCreate)
	err := c.ShouldBind(request)
	if err != nil {
		app.NewResponse(c).Error(errcode.ErrParams.WithCause(err))
		return
	}
	// 验证用户信息 Token 然后把UserID赋值上去 这里测试就直接赋值了
	request.UserId = 123453453
	svc := appservice.NewDemoAppSvc(c)
	reply, err := svc.CreateDemoOrder(request)
	if err != nil {
		app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
		return
	}
	app.NewResponse(c).Success(reply)
}
