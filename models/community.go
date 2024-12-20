package models

import (
	"ginchat/utils"

	"gorm.io/gorm"
)

type Community struct {
	gorm.Model
	Name    string //群名
	Ownerid uint   // 群主
	Number  int    // 群号
	Img     string //群图片
	Desc    string //群描述
}

func (table *Community) TableName() string {
	return "community"
}

// 创建群
func CreateCommunity(community Community) (bool, string) {
	if len(community.Name) == 0 {
		return false, "群名称不规范"
	}
	if community.Ownerid == 0 {
		return false, "群主不能为空"
	}
	if err := utils.DB.Create(&community).Error; err != nil {
		return false, "建群失败"
	}
	return true, "建群成功"
}

// 群列表
func Loadcommunity(id uint) []Community {
	communitys := make([]Community, 0)
	contacts := make([]Contact, 0)
	utils.DB.Where("ownerid = ? and type = 2", id).Find(&contacts)
	for _, v := range contacts {
		temp := Community{}
		utils.DB.Where("id = ?", v.Targetid).Find(&temp)
		communitys = append(communitys, temp)
	}
	return communitys
}

// 修改群名称
