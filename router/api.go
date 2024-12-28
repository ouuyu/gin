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
			selfRouter.GET("/reset", controller.ResetUserPassword)
		}
	}

	systemRouter := api.Group("/system")
	systemRouter.Use(middleware.AuthAdmin())
	{
		systemRouter.GET("/config", controller.GetSystemConfigs)
		systemRouter.PUT("/config", controller.UpdateSystemConfig)
		systemRouter.GET("/user/list", controller.GetUserList)
		systemRouter.POST("/user/update", controller.UpdateUser)
	}

	// 用户组管理路由
	groupController := &controller.GroupController{}
	groups := api.Group("/groups")
	{
		groups.POST("", groupController.CreateGroup)       // 创建用户组
		groups.GET("", groupController.GetAllGroups)       // 获取所有用户组
		groups.PUT("/:id", groupController.UpdateGroup)    // 更新用户组
		groups.DELETE("/:id", groupController.DeleteGroup) // 删除用户组
		groups.GET("/:id", groupController.GetGroup)       // 获取用户组信息
	}
}
