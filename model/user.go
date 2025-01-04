package model

import (
	"fmt"
	"main/common"
)

type User struct {
	ID               int     `json:"id" gorm:"type:int;primaryKey;autoIncrement"`
	Username         string  `json:"username" gorm:"unique;index" validate:"required,min=5,max=12"`
	Password         string  `json:"password" validate:"required"`
	Role             int     `json:"role" gorm:"default:0"`
	GroupId          int     `json:"group_id" gorm:"type:int;default:0"`
	Group            Group   `json:"group" gorm:"foreignKey:GroupId;references:ID"`
	Status           int     `json:"status" gorm:"type:int;default:1"` // 禁用用户为 2，正常用户为 1
	Token            string  `json:"token" gorm:"index"`
	Email            string  `json:"email" gorm:"index" validate:"max=50"`
	GitHubId         string  `json:"github_id" gorm:"column:github_id;index"`
	Balance          float64 `json:"balance" gorm:"type:decimal(10,2);default:0"` // 用户余额
	VerificationCode string  `json:"verification_code" gorm:"-:all"`
}

const cleanUserFields = "username, role, id, email, status, token, group_id, balance" // 不向前端发送敏感信息

func (u *User) Insert() error {
	var err error
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

func (u *User) Update() error {
	return DB.Model(u).Where("id = ?", u.ID).Updates(u).Error
}

func (u *User) UpdatePassword(password string) error {
	u.Password, _ = common.Password2Hash(password)
	return u.Update()
}

func (u *User) Delete(id int) error {
	return DB.Delete(u, id).Error
}

func (u *User) ValidateAndLogin() (*User, error) {
	var dbUser User
	if err := DB.Where("username = ?", u.Username).First(&dbUser).Error; err != nil {
		return nil, fmt.Errorf("用户不存在")
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
	query := DB
	if clean {
		query = query.Select(cleanUserFields)
	}
	err := query.Preload("Group").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserList(offset, pageSize int) ([]User, error) {
	var users []User
	err := DB.Preload("Group").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error
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

func (u *User) UpdateGroup(groupId int) error {
	u.GroupId = groupId
	return DB.Model(u).Update("group_id", groupId).Error
}

// UpdateBalance 更新用户余额
func (u *User) UpdateBalance(amount float64) error {
	// 开启事务
	tx := DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新余额
	if err := tx.Model(u).Update("balance", amount).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

// AddBalance 增加余额（充值）
func (u *User) AddBalance(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("充值金额必须大于0")
	}
	return u.UpdateBalance(u.Balance + amount)
}

// DeductBalance 扣除余额（消费）
func (u *User) DeductBalance(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("扣除金额必须大于0")
	}
	if u.Balance < amount {
		return fmt.Errorf("余额不足")
	}
	return u.UpdateBalance(u.Balance - amount)
}

// GetBalance 获取用户余额
func (u *User) GetBalance() float64 {
	return u.Balance
}
