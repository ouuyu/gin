package model

import (
	"main/common"
	"sync"
	"time"
)

type Config struct {
	Version    string `json:"version"`
	SystemName string `json:"system_name"`
	Footer     string `json:"footer"`
	HomePage   string `json:"home_page"`
	SQLitePath string `json:"sqlite_path"`

	RegisterEnabled          bool `json:"register_enabled"`
	PasswordRegisterEnabled  bool `json:"password_register_enabled"`
	EmailVerificationEnabled bool `json:"email_verification_enabled"`
	RecaptchaEnabled         bool `json:"recaptcha_enabled"`

	SMTPServer   string `json:"smtp_server"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUser     string `json:"smtp_user"`
	SMTPPassword string `json:"smtp_password"`
	SMTPFrom     string `json:"smtp_from"`

	RecaptchaSiteKey   string `json:"recaptcha_site_key"`
	RecaptchaSecretKey string `json:"recaptcha_secret_key"`

	RateLimitKeyExpirationDuration time.Duration `json:"rate_limit_key_expiration_duration"`
	GlobalApiRateLimitNum          int           `json:"global_api_rate_limit_num"`
	GlobalApiRateLimitDuration     int64         `json:"global_api_rate_limit_duration"`
	GlobalWebRateLimitNum          int           `json:"global_web_rate_limit_num"`
	GlobalWebRateLimitDuration     int64         `json:"global_web_rate_limit_duration"`
	UploadRateLimitNum             int           `json:"upload_rate_limit_num"`
	UploadRateLimitDuration        int64         `json:"upload_rate_limit_duration"`
	DownloadRateLimitNum           int           `json:"download_rate_limit_num"`
	DownloadRateLimitDuration      int64         `json:"download_rate_limit_duration"`
	CriticalRateLimitNum           int           `json:"critical_rate_limit_num"`
	CriticalRateLimitDuration      int64         `json:"critical_rate_limit_duration"`
}

var (
	config     *Config
	configLock sync.RWMutex
)

func GetConfig() *Config {
	configLock.RLock()
	if config != nil {
		defer configLock.RUnlock()
		return config
	}
	configLock.RUnlock()

	configLock.Lock()
	defer configLock.Unlock()

	var cfg Config
	if err := DB.First(&cfg).Error; err != nil {
		cfg = Config{
			Version:                        "0.0.1",
			SystemName:                     "管理系统",
			SQLitePath:                     "data.db",
			RegisterEnabled:                true,
			PasswordRegisterEnabled:        true,
			EmailVerificationEnabled:       false,
			RecaptchaEnabled:               false,
			SMTPPort:                       587,
			RateLimitKeyExpirationDuration: 20 * time.Minute,
			GlobalApiRateLimitNum:          60,
			GlobalApiRateLimitDuration:     3 * 60,
			GlobalWebRateLimitNum:          60,
			GlobalWebRateLimitDuration:     3 * 60,
			UploadRateLimitNum:             10,
			UploadRateLimitDuration:        60,
			DownloadRateLimitNum:           10,
			DownloadRateLimitDuration:      60,
			CriticalRateLimitNum:           20,
			CriticalRateLimitDuration:      20 * 60,
		}
		DB.Create(&cfg)
	}
	config = &cfg
	updateGlobalVars()
	return config
}

func SaveConfig(cfg Config) error {
	configLock.Lock()
	defer configLock.Unlock()

	// make the validator happy
	if err := DB.Where("1 = 1").Updates(&cfg).Error; err != nil {
		return err
	}
	config = &cfg
	updateGlobalVars()
	return nil
}

func updateGlobalVars() {
	common.Version = config.Version
	common.SystemName = config.SystemName
	common.Footer = config.Footer
	common.HomePage = config.HomePage
	common.SQLitePath = config.SQLitePath
	common.RegisterEnabled = config.RegisterEnabled
	common.PasswordRegisterEnabled = config.PasswordRegisterEnabled
	common.EmailVerificationEnabled = config.EmailVerificationEnabled
	common.RecaptchaEnabled = config.RecaptchaEnabled
	common.SMTPServer = config.SMTPServer
	common.SMTPPort = config.SMTPPort
	common.SMTPUser = config.SMTPUser
	common.SMTPPassword = config.SMTPPassword
	common.SMTPFrom = config.SMTPFrom
	common.RecaptchaSiteKey = config.RecaptchaSiteKey
	common.RecaptchaSecretKey = config.RecaptchaSecretKey
	common.RateLimitKeyExpirationDuration = config.RateLimitKeyExpirationDuration
	common.GlobalApiRateLimitNum = config.GlobalApiRateLimitNum
	common.GlobalApiRateLimitDuration = config.GlobalApiRateLimitDuration
	common.GlobalWebRateLimitNum = config.GlobalWebRateLimitNum
	common.GlobalWebRateLimitDuration = config.GlobalWebRateLimitDuration
	common.UploadRateLimitNum = config.UploadRateLimitNum
	common.UploadRateLimitDuration = config.UploadRateLimitDuration
	common.DownloadRateLimitNum = config.DownloadRateLimitNum
	common.DownloadRateLimitDuration = config.DownloadRateLimitDuration
	common.CriticalRateLimitNum = config.CriticalRateLimitNum
	common.CriticalRateLimitDuration = config.CriticalRateLimitDuration
}
