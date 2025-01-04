package controller

import (
	"encoding/json"
	"main/common"
	"main/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Register(c *gin.Context) {
	if !common.RegisterEnabled {
		common.Fail(c, http.StatusOK, "管理员关闭了新用户注册")
		return
	}
	if !common.PasswordRegisterEnabled {
		common.Fail(c, http.StatusOK, "管理员关闭了通过密码进行注册，请使用第三方账户验证的形式进行注册")
		return
	}
	var user model.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		common.Fail(c, http.StatusOK, "无效的参数")
		return
	}
	if err := common.Validate.Struct(&user); err != nil {
		common.Fail(c, http.StatusOK, "输入不合法: "+err.Error())
		return
	}
	if common.EmailVerificationEnabled {
		if user.Email == "" || user.VerificationCode == "" {
			common.Fail(c, http.StatusOK, "管理员开启了邮箱验证，请输入邮箱地址和验证码")
			return
		}
		if !common.VerifyCodeWithKey(user.Email, user.VerificationCode) {
			common.Fail(c, http.StatusOK, "验证码错误或已过期")
			return
		}
	}
	cleanUser := model.User{
		Username: user.Username,
		Password: user.Password,
	}
	if common.EmailVerificationEnabled {
		cleanUser.Email = user.Email
	}
	if err := cleanUser.Insert(); err != nil {
		errMessage := err.Error()
		if strings.Contains(errMessage, "UNIQUE constraint failed") {
			errMessage = "用户名已存在"
		}
		common.Fail(c, http.StatusOK, errMessage)
		return
	}
	common.Success(c, nil)
}

func Login(c *gin.Context) {
	var user model.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		common.Fail(c, http.StatusOK, "无效的参数")
		return
	}

	loginUser, err := user.ValidateAndLogin()
	if err != nil {
		common.Fail(c, http.StatusOK, "错误: "+err.Error())
		return
	}

	common.Success(c, gin.H{
		"token":    loginUser.Token,
		"username": loginUser.Username,
		"role":     loginUser.Role,
	})
}

func GenerateToken(c *gin.Context) {
	id := c.GetInt("id")
	user, err := model.GetUserById(id, false)
	if err != nil {
		common.Fail(c, http.StatusOK, err.Error())
		return
	}
	user.Token = uuid.New().String()
	user.Token = strings.Replace(user.Token, "-", "", -1)

	if model.DB.Where("token = ?", user.Token).First(user).RowsAffected != 0 {
		common.Fail(c, http.StatusOK, "请重试，系统生成的 UUID 重复？")
		return
	}

	if err := user.Update(); err != nil {
		common.Fail(c, http.StatusOK, err.Error())
		return
	}

	common.Success(c, user.Token)
}

func GetUserInfo(c *gin.Context) {
	id := c.GetInt("id")
	userid, _ := strconv.Atoi(c.DefaultQuery("userid", strconv.Itoa(id)))
	if userid != id && c.GetInt("role") < common.RoleAdmin {
		common.Fail(c, http.StatusOK, "您没有权限查看其他用户的信息")
		return
	}
	user, err := model.GetUserById(userid, true)
	if err != nil {
		common.Fail(c, http.StatusOK, err.Error())
		return
	}
	common.Success(c, user)
}

func GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	users, err := model.GetUserList(offset, pageSize)
	if err != nil {
		common.Fail(c, http.StatusOK, err.Error())
		return
	}

	total, err := model.GetUserCount()
	if err != nil {
		common.Fail(c, http.StatusOK, err.Error())
		return
	}

	common.Success(c, gin.H{
		"list":  users,
		"total": total,
	})
}

func ResetUserPassword(c *gin.Context) {
	id := c.GetInt("id")
	role := c.GetInt("role")
	userid, _ := strconv.Atoi(c.DefaultQuery("userid", strconv.Itoa(id)))
	password := c.DefaultQuery("password", "123456")

	if role == common.RoleUser {
		if userid != id {
			common.Fail(c, http.StatusOK, "您没有权限重置其他用户的密码")
			return
		}
	}

	user, err := model.GetUserById(userid, false)
	if err != nil {
		common.Fail(c, http.StatusOK, err.Error())
		return
	}

	if err := user.UpdatePassword(password); err != nil {
		common.Fail(c, http.StatusOK, err.Error())
		return
	}
	common.Success(c, nil)
}

func UpdateUser(c *gin.Context) {
	var user model.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		common.Fail(c, http.StatusOK, "无效的参数")
		return
	}

	user.Username = ""

	if err := user.Update(); err != nil {
		common.Fail(c, http.StatusOK, err.Error())
		return
	}

	common.Success(c, nil)
}

func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user := &model.User{ID: id}

	if err := user.Delete(id); err != nil {
		common.Fail(c, http.StatusOK, err.Error())
		return
	}

	common.Success(c, nil)
}

func CreateUser(c *gin.Context) {
	var user model.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		common.Fail(c, http.StatusOK, "无效的参数")
		return
	}

	user.Role = common.RoleUser
	// important: 密码哈希在 Insert 方法中进行

	if err := user.Insert(); err != nil {
		common.Fail(c, http.StatusOK, err.Error())
		return
	}

	common.Success(c, nil)
}
