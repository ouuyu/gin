package model

import (
	"fmt"
	"main/common"
)

type User struct {
	ID               int    `json:"id" gorm:"type:int;primaryKey"`
	Username         string `json:"username" gorm:"unique;index" validate:"max=12"`
	Password         string `json:"password" gorm:"not null;" validate:"min=8,max=30"`
	Role             int    `json:"role" gorm:"type:int;default:1"`   // root is 100, common user is 1
	Status           int    `json:"status" gorm:"type:int;default:1"` // enabled, disabled
	Token            string `json:"token" gorm:"index"`
	Email            string `json:"email" gorm:"index" validate:"max=50"`
	GitHubId         string `json:"github_id" gorm:"column:github_id;index"`
	VerificationCode string `json:"verification_code" gorm:"-:all"`
}

func (u *User) Insert() error {
	var err error
	// 验证未经过哈希的密码
	if err := common.Validate.Struct(u); err != nil {
		return err
	}
	if u.Password != "" {
		u.Password, err = common.Password2Hash(u.Password)
		if err != nil {
			return err
		}
		err = DB.Create(u).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (u *User) Update(clean bool) error {
	if clean {
		u.Password = ""
		u.Token = ""
	}
	return DB.Save(u).Error
}

func (u *User) ValidateAndLogin() (*User, error) {
	var dbUser User
	if err := DB.Where("username = ?", u.Username).First(&dbUser).Error; err != nil {
		return nil, err
	}

	if !common.ValidatePasswordAndHash(u.Password, dbUser.Password) {
		return nil, fmt.Errorf("密码错误")
	}

	if dbUser.Status != common.UserStatusEnabled {
		return nil, fmt.Errorf("账户已被禁用")
	}

	token, err := common.GenerateJwt(dbUser.Username, dbUser.Role, int(dbUser.ID))
	if err != nil {
		return nil, err
	}

	var cleanUser User
	cleanUser.Username = dbUser.Username
	cleanUser.Role = dbUser.Role
	cleanUser.Token = token

	return &cleanUser, nil
}

func GetUserById(id int, clean bool) (*User, error) {
	var user User
	if clean {
		DB.Select("username", "role", "id", "email", "status", "token").Where("id = ?", id).First(&user)
	} else {
		DB.Where("id = ?", id).First(&user)
	}
	return &user, nil
}

func GetUserList(offset, pageSize int) ([]User, error) {
	var users []User
	err := DB.Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
