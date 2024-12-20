package models

import (
	"fmt"
	"ginchat/utils"
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name         string
	Password     string
	Phone        string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Emails       string `valid:"email"`
	Avatar       string //头像
	Identity     string
	ClientIp     string
	ClientPort   string
	Salt         string
	LoginTime    time.Time
	HearBeatTime *time.Time
	LoginOutTime *time.Time `gorm:"column:login_out_time" json:"login_out_time"`
	IsLogOut     bool
	DeviceInfo   string
}

/*
		func (receiver ReceiverType) MethodName(parameters) (returnTypes) {
	    	// 方法体
		}
			receiver：接收者变量的名称
			ReceiverType：接收者的类型，可以是值类型或指针类型
*/
func (table *UserBasic) TableName() string { // 函数名前面的括号表示该方法的接收者
	return "user_basic"
}

func GetUserList() []*UserBasic {
	data := make([]*UserBasic, 10) // 定义一个空的 []*UserBasic 切片，这个切片会自动根据数据库中的结果扩展大小。
	utils.DB.Find(&data)
	return data
}

// 验证登录
func IsUser(name, password string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ? and password = ?", name, password).First(&user)
	// token 的加密
	str := fmt.Sprintf("%d", time.Now().Unix())
	temp := utils.MD5Encode(str)
	utils.DB.Model(&user).Where("id = ?", user.ID).Update("Identity", temp)
	return user
}

// 根据用户名查询用户
func GetUserByName(name string) UserBasic {
	user := UserBasic{}
	utils.DB.Where("name = ?", name).First(&user)
	return user
}

// 添加用户
func CreatUser(user UserBasic) *gorm.DB {
	return utils.DB.Create(&user)
}

// 删除用户
func DeleteUser(user UserBasic) *gorm.DB {
	return utils.DB.Delete(&user)
}

// 修改用户
func UpdateUser(user UserBasic) *gorm.DB {
	// 查找用户是否存在
	var exituser UserBasic
	err := utils.DB.First(&exituser, user.ID).Error
	if err != nil {
		return nil
	}

	return utils.DB.Model(&user).Updates(UserBasic{
		Name:     user.Name,
		Password: user.Password,
		Emails:   user.Emails,
		Phone:    user.Phone,
		Avatar:   user.Avatar,
	})
}
