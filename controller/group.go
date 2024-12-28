package controller

import (
	"main/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateGroup(c *gin.Context) {
	var group model.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "创建用户组失败: " + err.Error(),
		})
		return
	}

	if err := group.Insert(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建用户组失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "创建用户组成功",
		"data":    group,
	})
}

func UpdateGroup(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var group model.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "更新用户组失败: " + err.Error(),
		})
		return
	}

	group.ID = id
	if err := group.Update(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新用户组失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新用户组成功",
		"data":    group,
	})
}

// DeleteGroup 删除用户组
// @Summary 删除用户组
// @Tags 用户组管理
// @Produce json
// @Param id path int true "用户组ID"
// @Success 200 {object} gin.H
// @Router /groups/{id} [delete]
func DeleteGroup(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	group := &model.Group{ID: id}

	if err := group.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除用户组失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除用户组成功",
	})
}

// GetGroup 获取用户组信息
// @Summary 获取用户组信息
// @Tags 用户组管理
// @Produce json
// @Param id path int true "用户组ID"
// @Success 200 {object} gin.H
// @Router /groups/{id} [get]
func GetGroup(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	group, err := model.GetGroupById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "用户组不存在: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取用户组成功",
		"data":    group,
	})
}

// GetAllGroups 获取所有用户组
// @Summary 获取所有用户组
// @Tags 用户组管理
// @Produce json
// @Success 200 {object} gin.H
// @Router /groups [get]
func GetAllGroups(c *gin.Context) {
	groups, err := model.GetAllGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取用户组列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取用户组列表成功",
		"data":    groups,
	})
}
