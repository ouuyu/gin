package controller

import (
	"fmt"
	"main/common"
	"main/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCleanConfig(c *gin.Context) {
	config := model.GetConfig()

	// 用户前端拉取，不发送敏感信息
	cleanConfig := model.Config{
		Version:                  config.Version,
		SystemName:               config.SystemName,
		Footer:                   config.Footer,
		HomePage:                 config.HomePage,
		RegisterEnabled:          config.RegisterEnabled,
		PasswordRegisterEnabled:  config.PasswordRegisterEnabled,
		EmailVerificationEnabled: config.EmailVerificationEnabled,
		RecaptchaEnabled:         config.RecaptchaEnabled,
		RecaptchaSiteKey:         config.RecaptchaSiteKey,
	}

	common.Success(c, cleanConfig)
}

func GetSystemConfigs(c *gin.Context) {
	configs := model.GetConfig()

	common.Success(c, configs)
}

type UpdateConfigRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

func UpdateSystemConfig(c *gin.Context) {
	var req UpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Fail(c, http.StatusBadRequest, "无效的参数")
		return
	}

	config := model.GetConfig()

	switch req.Key {
	case "system_name":
		config.SystemName = req.Value
	case "footer":
		config.Footer = req.Value
	case "home_page":
		config.HomePage = req.Value
	case "register_enabled":
		if req.Value == "true" {
			config.RegisterEnabled = true
		} else {
			config.RegisterEnabled = false
		}
	case "password_register_enabled":
		if req.Value == "true" {
			config.PasswordRegisterEnabled = true
		} else {
			config.PasswordRegisterEnabled = false
		}
	case "email_verification_enabled":
		if req.Value == "true" {
			config.EmailVerificationEnabled = true
		} else {
			config.EmailVerificationEnabled = false
		}
	case "recaptcha_enabled":
		if req.Value == "true" {
			config.RecaptchaEnabled = true
		} else {
			config.RecaptchaEnabled = false
		}
	case "smtp_server":
		config.SMTPServer = req.Value
	case "smtp_port":
		port := 0
		fmt.Sscanf(req.Value, "%d", &port)
		if port > 0 {
			config.SMTPPort = port
		}
	case "smtp_user":
		config.SMTPUser = req.Value
	case "smtp_password":
		config.SMTPPassword = req.Value
	case "smtp_from":
		config.SMTPFrom = req.Value
	case "recaptcha_site_key":
		config.RecaptchaSiteKey = req.Value
	case "recaptcha_secret_key":
		config.RecaptchaSecretKey = req.Value
	case "easy_pay_url":
		config.EasyPayURL = req.Value
	case "easy_pay_pid":
		config.EasyPayPid = req.Value
	case "easy_pay_key":
		config.EasyPayKey = req.Value
	default:
		common.Fail(c, http.StatusBadRequest, "无效的配置项")
		return
	}

	if err := model.SaveConfig(*config); err != nil {
		common.Fail(c, http.StatusInternalServerError, "保存配置失败: "+err.Error())
		return
	}

	common.Success(c, config)
}
