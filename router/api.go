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
		systemRouter.POST("/user/create", controller.CreateUser)
		systemRouter.DELETE("/user/delete/:id", controller.DeleteUser)
	}

	groupRouter := api.Group("/group")
	groupRouter.Use(middleware.AuthAdmin())
	{
		groupRouter.POST("/create", controller.CreateGroup)
		groupRouter.GET("/list", controller.GetAllGroups)
		groupRouter.PUT("/update/:id", controller.UpdateGroup)
		groupRouter.DELETE("/delete/:id", controller.DeleteGroup)
		groupRouter.GET("/info/:id", controller.GetGroup)
	}

	// 订单相关路由
	orderRouter := api.Group("/order")
	orderRouter.Use(middleware.AuthUser())
	{
		orderRouter.GET("/list", controller.GetOrderList)
		orderRouter.GET("/query/:trade_no", controller.QueryOrder)
	}

	// 余额相关路由
	balanceRouter := api.Group("/balance")
	balanceRouter.Use(middleware.AuthUser())
	{
		balanceRouter.POST("/recharge", controller.Recharge)  // 充值
		balanceRouter.GET("/logs", controller.GetBalanceLogs) // 获取余额变动记录
		balanceRouter.GET("/info", controller.GetBalance)     // 获取余额信息
	}

}
