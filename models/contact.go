package models

import (
	"fmt"
	"ginchat/utils"

	"gorm.io/gorm"
)

type Contact struct {
	gorm.Model
	Ownerid  uint // 谁的关系信息
	Targetid uint // 对应的谁
	Type     int  // 对应的类型 1：好友 2：群
	Desc     string
}

func (table *Contact) TableName() string {
	return "contact"
}

// 查询好友
func SearchFriends(userId uint) []UserBasic {
	contacts := make([]Contact, 0)
	objIds := make([]uint64, 0)
	utils.DB.Where("ownerid = ? and type=1", userId).Find(&contacts)
	for _, v := range contacts {
		fmt.Println(v)
		objIds = append(objIds, uint64(v.Targetid))
	}
	users := make([]UserBasic, 0)
	utils.DB.Where("id in ?", objIds).Find(&users)
	return users
}

// 添加好友
func AddFriend(userId uint, targetid uint) (bool, string) {
	user := UserBasic{}
	if targetid != 0 {
		if targetid == userId {
			return false, "不能添加自己为好友"
		}
		var exitcontact Contact
		err := utils.DB.Where("ownerid = ? and targetid = ? and type=1", userId, targetid).First(&exitcontact).Error
		if err == nil {
			return false, "已添加过此好友"
		}
		utils.DB.Where("id = ?", targetid).Find(&user)
		if user.Salt != "" {
			/*
				事务（Transaction）通常用于数据库操作，它保证了一组数据库操作要么全部成功，
				要么全部失败，以确保数据的一致性和可靠性。
				特性：原子性：事务中的所有操作要么都成功，要么都失败，不会有中间状态。
					一致性：事务执行前后的数据库必须符合所有的业务规则和约束条件。
					隔离性：一个事务的操作对其他事务是隔离的，不会被其他事务干扰，避免了并发问题。
					持久性：一旦事务提交，对数据库的更改是永久的，不会丢失，即使发生系统崩溃。
			*/
			// 使用事务对象执行操作
			session := utils.DB.Begin()
			// 事务一旦开始，无论什么异常都会回滚事务
			defer func() {
				// recover()是一个内置函数，它用于捕获程序中的 panic
				if r := recover(); r != nil {
					session.Rollback()
				}
			}()
			contact1 := Contact{
				Ownerid:  userId,
				Targetid: targetid,
				Type:     1,
			}
			if err := session.Create(&contact1).Error; err != nil {
				return false, "添加好友失败"
			}
			contact2 := Contact{
				Ownerid:  targetid,
				Targetid: userId,
				Type:     1,
			}
			if err := session.Create(&contact2).Error; err != nil {
				return false, "添加好友失败"
			}
			// 提交事务
			if err := session.Commit().Error; err != nil {
				session.Rollback()
				return false, "添加好友失败"
			}
			return true, "添加好友成功"
		}
		return false, "没有找到这个用户"
	}
	return false, "不存在这个用户"
}

// 加入群聊
func JoinGroup(userid uint, number int) (bool, string) {
	community := Community{}
	contact := Contact{}
	utils.DB.Where("number = ?", number).Find(&community)
	if community.Name == "" {
		return false, "没有找到这个群聊"
	}
	utils.DB.Where("ownerid = ? and targetid = ? and type = 2", userid, community.ID).First(&contact)
	if !contact.CreatedAt.IsZero() {
		return false, "你已经加过群，不需要重复加入"
	}
	contact.Ownerid = userid
	contact.Targetid = community.ID
	contact.Type = 2
	utils.DB.Create(&contact)

	return true, "加群成功"
}

func SearchUserByGroupId(communityId uint) []uint {
	contacts := make([]Contact, 0)
	utils.DB.Where("targetid = ? and type = 2", communityId).Find(&contacts)
	objIds := make([]uint, 0)
	for _, v := range contacts {
		objIds = append(objIds, v.Ownerid)
	}
	return objIds
}
