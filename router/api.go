package router

import (
	"main/controller"

	"github.com/gin-gonic/gin"
)

func SetApiRouter(router *gin.Engine) {
	api := router.Group("/api")
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	userRouter := api.Group("/user")
	{
		userRouter.POST("/register", controller.Register)
	}
}
