package main

import (
	"ginchat/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/ginchat?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 迁移schema
	// db.AutoMigrate(&models.UserBasic{})
	db.AutoMigrate(&models.Message{})
	// db.AutoMigrate(&models.Contact{})
	// db.AutoMigrate(&models.Group{})
	// db.AutoMigrate(&models.Community{})

	// // 创建
	// user := &models.UserBasic{}
	// user.Name = "刘"
	// user.LoginTime = time.Now().Truncate(time.Second)
	// db.Create(user)

	// // 查找
	// fmt.Println(db.First(user, 1))

	// db.Model(user).Update("Password", "1234")

}
