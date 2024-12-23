package router

import (
	"main/controller"
	"main/middleware"

	"github.com/gin-gonic/gin"
)

func SetApiRouter(router *gin.Engine) {
	api := router.Group("/api")
	api.GET("/config", controller.GetCleanConfig)

	userRouter := api.Group("/user")
	{
		userRouter.POST("/register", controller.Register)
		userRouter.POST("/login", controller.Login)

		selfRouter := userRouter.Group("/")
		selfRouter.Use(middleware.AuthUser())
		{
			selfRouter.GET("/info", controller.GetUserInfo)
			selfRouter.GET("/token", controller.GenerateToken)
		}
	}

	systemRouter := api.Group("/system")
	systemRouter.Use(middleware.AuthAdmin())
	{
		systemRouter.GET("/config", controller.GetSystemConfigs)
		systemRouter.PUT("/config", controller.UpdateSystemConfig)
		systemRouter.GET("/users", controller.GetUserList)
	}
}
