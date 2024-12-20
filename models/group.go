package models

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name    string
	Ownerid uint
	Icon    string
	Type    int
	Desc    string
}

func (table *Group) TableName() string {
	return "group_basic"
}
