package controller

import (
	"main/common"
	"main/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateGroup(c *gin.Context) {
	var group model.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		common.Fail(c, http.StatusBadRequest, "创建用户组失败: "+err.Error())
		return
	}

	if err := group.Insert(); err != nil {
		common.Fail(c, http.StatusInternalServerError, "创建用户组失败: "+err.Error())
		return
	}

	common.Success(c, group)
}

func UpdateGroup(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var group model.Group
	if err := c.ShouldBindJSON(&group); err != nil {
		common.Fail(c, http.StatusBadRequest, "更新用户组失败: "+err.Error())
		return
	}

	group.ID = id
	if err := group.Update(); err != nil {
		common.Fail(c, http.StatusInternalServerError, "更新用户组失败: "+err.Error())
		return
	}

	common.Success(c, group)
}

func DeleteGroup(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	group := &model.Group{ID: id}

	if err := group.Delete(false); err != nil {
		common.Fail(c, http.StatusInternalServerError, "删除用户组失败: "+err.Error())
		return
	}

	common.Success(c, nil)
}

func GetGroup(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	group, err := model.GetGroupById(id)
	if err != nil {
		common.Fail(c, http.StatusNotFound, "用户组不存在: "+err.Error())
		return
	}

	common.Success(c, group)
}

func GetAllGroups(c *gin.Context) {
	groups, err := model.GetAllGroups()
	if err != nil {
		common.Fail(c, http.StatusInternalServerError, "获取用户组列表失败: "+err.Error())
		return
	}

	common.Success(c, groups)
}
