package common

var (
	Version = "0.0.1"
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
	EmailVerificationPurpose = "email_verification"
)
