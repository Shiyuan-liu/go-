package service

import (
	"fmt"
	"ginchat/models"
	"ginchat/utils"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// GetUserList
// @summary 所有用户
// @Tags 用户模块
// @Success 200 {string} json{"code","message"}
// @Router /user/GetUserList [get]
func GetUserList(c *gin.Context) {
	data := make([]*models.UserBasic, 10)
	data = models.GetUserList()
	c.JSON(200, gin.H{
		"code":    0,
		"message": "查询成功！",
		"data":    data,
	})
}

// CreateUser
// @summary 添加用户
// @Tags 用户模块
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/CreateUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	// user.Name = c.Query("name")
	// password := c.Query("password")
	// repassword := c.Query("repassword")
	user.Name = c.Request.FormValue("name")
	password := c.Request.FormValue("password")
	repassword := c.Request.FormValue("repassword")

	salt := fmt.Sprintf("%06d", rand.Int31())

	if user.Name == "" || password == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "用户名或密码不能为空",
		})
		return
	}
	// 判断用户名是否存在
	data := models.GetUserByName(user.Name)
	if data.Name != "" {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "用户名已经存在",
		})
		return
	}
	user.LoginTime = time.Now().Truncate(time.Second) // 设置为当前时间
	if password != repassword {
		c.JSON(-1, gin.H{
			"code":    -1,
			"message": "两次密码不一致！",
		})
	} else {
		user.Password = utils.MakePassword(password, salt)
		user.Salt = salt
		models.CreatUser(user)
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "新增用户成功！",
		})
	}
}

// DeleteUser
// @summary 删除用户
// @Tags 用户模块
// @param id query string false "id"
// @Success 200 {string} json{"code","message"}
// @Router /user/DeleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	models.DeleteUser(user)
	c.JSON(200, gin.H{
		"code":    0,
		"message": "删除用户成功！",
	})
}

// UpdateUser
// @summary 修改用户
// @Tags 用户模块
// @param id formData string false "id"
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param emails formData string false "邮箱"
// @param phone formData string false "电话"
// @Success 200 {string} json{"code","message"}
// @Router /user/UpdateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Emails = c.PostForm("emails")
	user.Phone = c.PostForm("phone")
	user.Avatar = c.PostForm("icon")
	_, err := govalidator.ValidateStruct(user)
	if err == nil {
		models.UpdateUser(user)
		c.JSON(200, gin.H{
			"code":    0,
			"message": "修改用户成功！",
		})
	} else {
		fmt.Println("error：", err)
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "修改用户失败！",
		})
	}
}

// IsUser
// @summary 登录验证
// @Tags 登录
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/Isuser [post]
func IsUser(c *gin.Context) {

	data := models.UserBasic{}

	// name := c.Query("name")
	// password := c.Query("password")
	name := c.Request.FormValue("name")
	password := c.Request.FormValue("password")

	// 判断用户名是否存在
	user := models.GetUserByName(name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "用户不存在",
		})
		return
	}
	// 解密
	flag := utils.ValidPassword(password, user.Salt, user.Password)
	if !flag {
		c.JSON(200, gin.H{
			"code":    -1,
			"message": "密码不正确",
		})
		return
	}

	pwd := utils.MakePassword(password, user.Salt)
	data = models.IsUser(name, pwd)

	c.JSON(200, gin.H{
		"code":    0, // 0成功   -1失败
		"message": "登陆成功！",
		"data":    data,
	})

}

// 根据Id查找用户
func FindById(c *gin.Context) {
	userid, err := strconv.Atoi(c.Request.FormValue("userId"))
	if err != nil {
		fmt.Println(err)
		return
	}
	user := models.UserBasic{}
	utils.DB.Where("id = ?", userid).First(&user)
	data := models.GetUserByName(user.Name)
	utils.RespOk(c.Writer, "ok", data)
}

// 防止跨域站点伪造请求
var upGrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SendMsg(c *gin.Context) {
	// 通过websocket.Upgrader将HTTP连接升级成WebSocket连接
	ws, err := upGrade.Upgrade(c.Writer, c.Request, nil) // ws是通过WebSocket连接到客户端的一个对象。代表了一个WebSocket连接，所有后续的消息交互都会通过这个对象进行。
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func(ws *websocket.Conn) {
		err = ws.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(ws)
	MsgHandler(ws, c)
}

func MsgHandler(ws *websocket.Conn, c *gin.Context) {
	for {
		msg, err := utils.Subscribe(c, utils.PublishKey) // 通过Subscribe方法订阅redis获取消息
		if err != nil {
			fmt.Println(err)
			break // 发送错误退出循环，关闭连接
		}
		// 获取当前时间并格式化
		tm := time.Now().Format("2006-01-02 15:05:03")
		m := fmt.Sprintf("[ws][%s]:%s", tm, msg)
		// WebSocket通过WriteMessage方法向客户端发送消息
		err = ws.WriteMessage(1, []byte(m))
		if err != nil {
			fmt.Println(err)
			break // 如果写入消息失败，退出循环，关闭连接
		}
	}
}

// 与某人聊天发送消息
func SendUserMsg(c *gin.Context) {
	models.Chat(c.Writer, c.Request)
}

// 查询好友列表
func SearchFriends(c *gin.Context) {
	id, _ := strconv.Atoi(c.Request.FormValue("userId"))
	users := models.SearchFriends(uint(id))
	// c.JSON(http.StatusOK, gin.H{
	// 	"code":    0,
	// 	"message": "查询好友列表成功",
	// 	"data":    users,
	// })
	utils.RespOkList(c.Writer, len(users), users)
}

// redis发送消息
func RedisMsg(c *gin.Context) {
	userIdA, _ := strconv.Atoi(c.PostForm("userIdA"))
	userIdB, _ := strconv.Atoi(c.PostForm("userIdB"))
	start, _ := strconv.Atoi(c.PostForm("start"))
	end, _ := strconv.Atoi(c.PostForm("end"))
	isRev, _ := strconv.ParseBool(c.PostForm("isRev"))
	str := models.RedisMsg(userIdA, userIdB, start, end, isRev)
	utils.RespOkList(c.Writer, str, "ok")
}

// 添加好友
func AddFriend(c *gin.Context) {
	userId, err := strconv.Atoi(c.Request.FormValue("userId"))
	if err != nil {
		utils.RespFail(c.Writer, "Invalid userId")
		return
	}
	name := c.Request.FormValue("targetName")
	user := models.UserBasic{}
	utils.DB.Where("name = ?", name).Find(&user)
	targetId := user.ID
	k, str := models.AddFriend(uint(userId), targetId)
	if k {
		utils.RespOk(c.Writer, str, k)
	} else {
		utils.RespFail(c.Writer, str)
	}
	// targetId, err := strconv.Atoi(c.Request.FormValue("targetId"))
	// if err != nil {
	// 	utils.RespFail(c.Writer, "Invalid targetId")
	// 	return
	// }
}

// 建群
func CreateCommunity(c *gin.Context) {
	name := c.Request.FormValue("name")
	ownerid, err := strconv.Atoi(c.Request.FormValue("ownerid"))
	icon := c.Request.FormValue("icon")
	desc := c.Request.FormValue("desc")
	if err != nil {
		return
	}
	// 随机生成群号
	rand.Seed(time.Now().UnixNano())
	numDigits := rand.Intn(2) + 8
	lowerBound := 1
	for i := 1; i < numDigits; i++ {
		lowerBound *= 10
	}
	upperBound := lowerBound * 10
	number := rand.Intn(upperBound-lowerBound) + lowerBound
	community := models.Community{
		Name:    name,
		Ownerid: uint(ownerid),
		Number:  number,
		Desc:    desc,
		Img:     icon,
	}
	k, str := models.CreateCommunity(community)
	if k {
		bo, _ := models.JoinGroup(uint(ownerid), number)
		if bo {
			utils.RespOk(c.Writer, str, k)
		} else {
			utils.RespFail(c.Writer, "由于不能添加群主，建群失败")
		}
	} else {
		utils.RespFail(c.Writer, str)
	}
}

// 查询群列表
func Loadcommunity(c *gin.Context) {
	id, err := strconv.Atoi(c.Request.FormValue("ownerId"))
	if err != nil {
		return
	}
	communitys := models.Loadcommunity(uint(id))
	if len(communitys) != 0 {
		utils.RespOkList(c.Writer, len(communitys), communitys)
	} else {
		utils.RespFail(c.Writer, "查询失败")
	}

}

// 加群
func JoinGroup(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Request.FormValue("userId"))
	number, _ := strconv.Atoi(c.Request.FormValue("comId"))
	k, str := models.JoinGroup(uint(userId), number)
	if k {
		utils.RespOk(c.Writer, str, k)
	} else {
		utils.RespFail(c.Writer, str)
	}
}
