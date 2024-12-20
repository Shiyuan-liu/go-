package service

import (
	"fmt"
	"ginchat/models"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetIndex
// @Tags 首页
// @Success 200 {string} welcome
// @Router /index [get]
func GetIndex(c *gin.Context) {
	te, err := template.ParseFiles("index.html", "views/chat/head.html")
	if err != nil {
		panic(err)
	}
	te.Execute(c.Writer, "index")
	// c.JSON(200, gin.H{
	// 	"message": "welcome",
	// })
}

func Register(c *gin.Context) {
	te, err := template.ParseFiles("views/user/register.html")
	if err != nil {
		panic(err)
	}
	te.Execute(c.Writer, "register")
}

func ToChat(c *gin.Context) {
	te, err := template.ParseFiles("views/chat/index.html",
		"views/chat/head.html",
		"views/chat/foot.html",
		"views/chat/tabmenu.html",
		"views/chat/concat.html",
		"views/chat/group.html",
		"views/chat/profile.html",
		"views/chat/main.html",
		"views/chat/userinfo.html",
		"views/chat/createcom.html",
	)
	if err != nil {
		panic(err)
	}
	userId, _ := strconv.Atoi(c.Query("userId"))
	token := c.Query("token")

	user := models.UserBasic{}
	user.ID = uint(userId)
	user.Identity = token

	// fmt.Println("chat---------", user)
	if err := te.Execute(c.Writer, user); err != nil {
		fmt.Println("error-----", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "template rendering", "error": err.Error()})
	}
}

func Chat(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}
