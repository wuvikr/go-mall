package main

import (
	"errors"
	"go-mall/common/app"
	"go-mall/common/errcode"
	"go-mall/common/logger"
	"go-mall/common/middleware"
	"go-mall/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.New()

	g.Use(middleware.LogAccess(), middleware.Recovery(), middleware.StartTrace())

	g.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	g.GET("/config-read", func(c *gin.Context) {
		database := config.Database
		log := config.App.Log
		c.JSON(http.StatusOK, gin.H{
			"type":        database.Type,
			"dsn":         database.DSN,
			"maxopen":     database.MaxOpenConn,
			"maxidle":     database.MaxIdleConn,
			"maxlifetime": database.MaxLifeTime,
			"log.path":    log.FilePath,
			"log.maxsize": log.FileMaxSize,
			"log.maxage":  log.BackUpFileMaxAge,
		})
	})

	g.GET("/logger-test", func(c *gin.Context) {
		logger.New(c).Info("logger-test", "key", "value", "key2", 2)

		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	g.GET("/access-log-test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	g.GET("/panic-log-test", func(c *gin.Context) {
		var a map[string]string
		a["k"] = "v"
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"data":   a,
		})
	})

	g.GET("/customized-error-test", func(c *gin.Context) {

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

	})

	g.GET("/response-object-test", func(c *gin.Context) {
		data := map[string]string{
			"key":  "value",
			"key2": "value2",
		}
		app.NewResponse(c).Success(data)
	})

	g.GET("/response-error-test", func(c *gin.Context) {
		baseErr := errors.New("a dao error")
		err := errcode.Wrap("包装错误", baseErr)
		app.NewResponse(c).Error(errcode.ErrServer.WithCause(err))
	})

	g.Run("127.0.0.1:8080")
}
