package common

import "time"

var (
	Version    = "0.0.1"
	SystemName = "管理系统"
	Footer     = ""
	HomePage   = ""
)

var (
	SQLitePath = "data.db"
)

const (
	RoleRootUser = 100
	RoleAdmin    = 10
	RoleUser     = 1
	RoleGuest    = 0
)

const (
	UserStatusEnabled  = 1
	UserStatusDisabled = 2
)

var (
	RegisterEnabled          = true
	PasswordRegisterEnabled  = true
	EmailVerificationEnabled = false
	RecaptchaEnabled         = false
)

var (
	SMTPServer   = ""
	SMTPPort     = 587
	SMTPUser     = ""
	SMTPPassword = ""
	SMTPFrom     = ""
)

var (
	RecaptchaSiteKey   = ""
	RecaptchaSecretKey = ""
)

// API 限流 单位: 秒
// 不可超过限流键的过期时间
var (
	RateLimitKeyExpirationDuration = 20 * time.Minute

	GlobalApiRateLimitNum            = 60
	GlobalApiRateLimitDuration int64 = 3 * 60

	GlobalWebRateLimitNum            = 60
	GlobalWebRateLimitDuration int64 = 3 * 60

	UploadRateLimitNum            = 10
	UploadRateLimitDuration int64 = 60

	DownloadRateLimitNum            = 10
	DownloadRateLimitDuration int64 = 60

	CriticalRateLimitNum            = 20
	CriticalRateLimitDuration int64 = 20 * 60
)
