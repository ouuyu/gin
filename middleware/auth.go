package middleware

import (
	"fmt"
	"main/common"
	"main/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func authHelper(c *gin.Context, minRole int) {
	token := c.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "未提供认证token",
		})
		c.Abort()
		return
	}
	claims, err := common.ValidateJwt(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": fmt.Sprintf("无效的token: %v", err),
		})
		c.Abort()
		return
	}
	user, err := model.GetUserById(claims.ID, false)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "用户不存在",
		})
		c.Abort()
		return
	}
	if user.Role < minRole {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "权限不足",
		})
		c.Abort()
		return
	}
	if user.Status != common.UserStatusEnabled {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "用户已被禁用",
		})
		c.Abort()
		return
	}
	c.Set("username", user.Username)
	c.Set("role", user.Role)
	c.Set("id", user.ID)
	c.Next()
}

func AuthAdmin() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHelper(c, common.RoleAdmin)
	}
}

func AuthUser() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHelper(c, common.RoleUser)
	}
}
