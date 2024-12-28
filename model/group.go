package model

import (
	"errors"
	"main/common"
)

type Group struct {
	ID        int    `json:"id" gorm:"type:int;primaryKey"`
	Name      string `json:"name" gorm:"unique;index" validate:"max=12"`
	UserCount int    `json:"user_count" gorm:"-"`
}

const cleanGroupFields = "id, name, user_count"

func (g *Group) Insert() error {
	if err := common.Validate.Struct(g); err != nil {
		return err
	}
	return DB.Create(g).Error
}

func (g *Group) Update() error {
	if err := common.Validate.Struct(g); err != nil {
		return err
	}
	return DB.Save(g).Error
}

func (g *Group) Delete(force bool) error {
	users, err := GetUsersByGroupId(g.ID)
	if err != nil {
		return err
	}
	if len(users) > 0 && !force {
		return errors.New("用户组内有用户，请强制删除")
	}
	// 强制删除，将该组的所有用户的组ID设为0
	if err := DB.Model(&User{}).Where("group_id = ?", g.ID).Update("group_id", 0).Error; err != nil {
		return err
	}
	return DB.Delete(g).Error
}

func GetUsersByGroupId(id int) ([]User, error) {
	var users []User
	err := DB.Where("group_id = ?", id).Find(&users).Error
	return users, err
}

func GetGroupById(id int) (*Group, error) {
	var group Group
	if err := DB.First(&group, id).Select(cleanGroupFields).Error; err != nil {
		return nil, err
	}

	count, err := group.GetUserCount()
	if err != nil {
		return nil, err
	}
	group.UserCount = count

	return &group, nil
}

func GetAllGroups() ([]Group, error) {
	var groups []Group
	if err := DB.Find(&groups).Error; err != nil {
		return nil, err
	}

	for i := range groups {
		count, err := groups[i].GetUserCount()
		if err != nil {
			return nil, err
		}
		groups[i].UserCount = count
	}

	return groups, nil
}

func (g *Group) GetUserCount() (int, error) {
	var count int64
	err := DB.Model(&User{}).Where("group_id = ?", g.ID).Count(&count).Error
	return int(count), err
}
