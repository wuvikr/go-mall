package main

import (
	"go-mall/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	g := gin.New()

	g.Use(gin.Logger(), gin.Recovery())

	g.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	g.GET("/config-read", func(c *gin.Context) {
		database := config.Database
		c.JSON(http.StatusOK, gin.H{
			"type":        database.Type,
			"dsn":         database.DSN,
			"maxopen":     database.MaxOpenConn,
			"maxidle":     database.MaxIdleConn,
			"maxlifetime": database.MaxLifeTime,
		})
	})

	g.Run(":8080")
}
