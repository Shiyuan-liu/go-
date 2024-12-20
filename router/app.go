package router

import (
	"ginchat/docs"
	"ginchat/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	/**
	swag的主要作用：
		生成 API 文档：swag 通过解析代码中的注释，自动生成符合 Swagger 规范的 API 文档，通常以 JSON 或 YAML 格式输出。
		文档更新：随着代码的更改，swag 可以自动更新 API 文档，确保文档与实际代码保持一致。
		提高开发效率：使用 swag 可以大幅度减少手动编写和维护 API 文档的工作量。
	*/)

func Router() *gin.Engine {
	r := gin.Default()

	//docs.SwaggerInfo.BasePath 是一个用于指定 API 基础路径的字段，可以帮助生成的 Swagger 文档正确显示 API 的完整路径。
	docs.SwaggerInfo.BasePath = ""
	//添加swagger路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	//静态资源
	r.Static("/asset", "asset/")
	r.LoadHTMLGlob("views/**/*")

	// 首页
	r.GET("/", service.GetIndex)
	r.GET("/index", service.GetIndex)
	r.GET("/register", service.Register)
	r.GET("/ToChat", service.ToChat)
	r.GET("/Chat", service.Chat)

	// 查询好友
	r.POST("/SearchFriends", service.SearchFriends)
	// 用户模块
	r.POST("/user/GetUserList", service.GetUserList)
	r.POST("/user/CreateUser", service.CreateUser)
	r.POST("/user/DeleteUser", service.DeleteUser)
	r.POST("/user/UpdateUser", service.UpdateUser)
	r.POST("/user/Isuser", service.IsUser)
	r.POST("/user/Find", service.FindById)

	// websocket发送消息
	r.GET("/user/SendMsg", service.SendMsg)
	r.GET("/user/SendUserMsg", service.SendUserMsg)

	r.POST("/user/RedisMsg", service.RedisMsg)

	// 上传图片/音频等文件
	r.POST("/attach/upload", service.Upload)
	// 添加好友
	r.POST("/contact/AddFriend", service.AddFriend)
	// 建群
	r.POST("/contact/CreateCommunity", service.CreateCommunity)
	r.POST("/contact/Loadcommunity", service.Loadcommunity)
	// 加群
	r.POST("/contact/JoinGroup", service.JoinGroup)

	return r
}
