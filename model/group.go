package model

import (
	"main/common"
)

type Group struct {
	ID   int    `json:"id" gorm:"type:int;primaryKey"`
	Name string `json:"name" gorm:"unique;index" validate:"min=1,max=12"`
}

const cleanGroupFields = "id, name"

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

func (g *Group) Delete() error {
	// 删除组之前，将该组的所有用户的组ID设为0
	if err := DB.Model(&User{}).Where("group_id = ?", g.ID).Update("group_id", 0).Error; err != nil {
		return err
	}
	return DB.Delete(g).Error
}

func GetGroupById(id int) (*Group, error) {
	var group Group
	err := DB.Where("id = ?", id).Select(cleanGroupFields).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func GetAllGroups() ([]Group, error) {
	var groups []Group
	err := DB.Select(cleanGroupFields).Find(&groups).Error
	return groups, err
}
