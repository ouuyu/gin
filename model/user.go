package model

import (
	"fmt"
	"main/common"
)

type User struct {
	ID               int    `json:"id" gorm:"type:int;primaryKey"`
	Username         string `json:"username" gorm:"unique;index" validate:"required,min=5,max=12"`
	Password         string `json:"password" validate:"required,min=6,max=20"`
	Role             int    `json:"role" gorm:"default:0"`
	GroupId          int    `json:"group_id" gorm:"default:0"`        // 用户所属的组ID
	GroupName        string `json:"group_name" gorm:"-"`              // 组名称，不存储在数据库中
	Status           int    `json:"status" gorm:"type:int;default:1"` // enabled, disabled
	Token            string `json:"token" gorm:"index"`
	Email            string `json:"email" gorm:"index" validate:"max=50"`
	GitHubId         string `json:"github_id" gorm:"column:github_id;index"`
	VerificationCode string `json:"verification_code" gorm:"-:all"`
}

const cleanUserFields = "username, role, id, email, status, token" // 不向前端发送敏感信息

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

func (u *User) Update(id int) error {
	return DB.Model(u).Where("id = ?", id).Updates(u).Error
}

func (u *User) UpdatePassword(password string) error {
	u.Password, _ = common.Password2Hash(password)
	return u.Update(u.ID)
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
	var err error
	if clean {
		err = DB.Select(cleanUserFields).Where("id = ?", id).First(&user).Error
	} else {
		err = DB.Where("id = ?", id).First(&user).Error
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserList(offset, pageSize int) ([]User, error) {
	var users []User
	err := DB.Offset(offset).Limit(pageSize).Select(cleanUserFields).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetUserCount() (int64, error) {
	var count int64
	err := DB.Model(&User{}).Count(&count).Error
	return count, err
}

// UpdateGroup 更新用户的组信息
func (u *User) UpdateGroup(groupId int) error {
	u.GroupId = groupId
	return DB.Model(u).Update("group_id", groupId).Error
}

// GetUsersByGroupId 获取同组的所有用户
func GetUsersByGroupId(groupId int) ([]User, error) {
	var users []User
	err := DB.Where("group_id = ?", groupId).Find(&users).Error
	return users, err
}
