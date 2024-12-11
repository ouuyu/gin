package model

type User struct {
	ID               uint   `json:"id" gorm:"primaryKey"`
	Username         string `json:"username" gorm:"unique;index" validate:"max=12"`
	Password         string `json:"password" gorm:"not null;" validate:"min=8,max=20"`
	Role             int    `json:"role" gorm:"type:int;default:1"`   // admin, common user
	Status           int    `json:"status" gorm:"type:int;default:1"` // enabled, disabled
	Token            string `json:"token" gorm:"index"`
	Email            string `json:"email" gorm:"index" validate:"max=50"`
	GitHubId         string `json:"github_id" gorm:"column:github_id;index"`
	VerificationCode string `json:"verification_code" gorm:"-:all"`
}
