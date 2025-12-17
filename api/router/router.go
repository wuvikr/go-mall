package router

import (
	"go-mall/common/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(engine *gin.Engine) {
	engine.Use(middleware.StartTrace(),middleware.LogAccess(),middleware.Recovery())
	routeGroup := engine.Group("")

	registerBuildingRoutes(routeGroup)
}