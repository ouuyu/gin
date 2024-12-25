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
		c.JSON(http.StatusOK, gin.H{
			"message": "管理员关闭了新用户注册",
			"success": false,
		})
		return
	}
	if !common.PasswordRegisterEnabled {
		c.JSON(http.StatusOK, gin.H{
			"message": "管理员关闭了通过密码进行注册，请使用第三方账户验证的形式进行注册",
			"success": false,
		})
		return
	}
	var user model.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}
	if err := common.Validate.Struct(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "输入不合法 " + err.Error(),
		})
		return
	}
	if common.EmailVerificationEnabled {
		if user.Email == "" || user.VerificationCode == "" {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "管理员开启了邮箱验证，请输入邮箱地址和验证码",
			})
			return
		}
		if !common.VerifyCodeWithKey(user.Email, user.VerificationCode) {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "验证码错误或已过期",
			})
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
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": errMessage,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "注册成功",
	})
}

func Login(c *gin.Context) {
	var user model.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}

	loginUser, err := user.ValidateAndLogin()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "用户名或密码错误",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
		"data": gin.H{
			"token":    loginUser.Token,
			"username": loginUser.Username,
			"role":     loginUser.Role,
		},
	})
}

func GenerateToken(c *gin.Context) {
	id := c.GetInt("id")
	user, err := model.GetUserById(id, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	user.Token = uuid.New().String()
	user.Token = strings.Replace(user.Token, "-", "", -1)

	if model.DB.Where("token = ?", user.Token).First(user).RowsAffected != 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "请重试，系统生成的 UUID 重复？",
		})
		return
	}

	if err := user.Update(user.ID); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "成功生成新的令牌",
		"data":    user.Token,
	})
}

func GetUserInfo(c *gin.Context) {
	id := c.GetInt("id")
	user, err := model.GetUserById(id, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	var cleanUser model.User
	cleanUser.Username = user.Username
	cleanUser.Role = user.Role
	cleanUser.ID = user.ID
	cleanUser.Email = user.Email
	cleanUser.Token = user.Token
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "成功获取用户信息",
		"data":    cleanUser,
	})
}

func GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	offset := (page - 1) * pageSize

	users, err := model.GetUserList(offset, pageSize)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	total, err := model.GetUserCount()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "成功获取用户列表",
		"data": gin.H{
			"list":  users,
			"total": total,
		},
	})
}

func ResetUserPassword(c *gin.Context) {
	id := c.GetInt("id")
	role := c.GetInt("role")
	userid, _ := strconv.Atoi(c.DefaultQuery("userid", strconv.Itoa(id)))
	password := c.DefaultQuery("password", "123456")

	if role == common.RoleUser {
		if userid != id {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "您没有权限重置其他用户的密码",
			})
			return
		}
	}

	user, err := model.GetUserById(userid, false)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if err := user.UpdatePassword(password); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "成功重置用户密码",
	})
}

func UpdateUser(c *gin.Context) {
	var user model.User
	err := json.NewDecoder(c.Request.Body).Decode(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "无效的参数",
		})
		return
	}

	user.Username = ""
	user.Password = ""
	user.Token = ""

	if err := user.Update(user.ID); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "成功更新用户信息",
	})
}
