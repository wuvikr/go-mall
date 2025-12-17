package main

import (
	"go-mall/api/router"
	"go-mall/common/enum"
	"go-mall/config"

	"github.com/gin-gonic/gin"
)

func main() {
	if config.App.Env == enum.ModeProd {
		gin.SetMode(gin.ReleaseMode)
	}
	g := gin.New()
	router.RegisterRoutes(g)

	g.Run(":8080")
}
