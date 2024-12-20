package models

import (
	"context"
	"encoding/json"
	"fmt"
	"ginchat/utils"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Userid     uint   //发送者
	Targetid   uint   //接收者
	Type       int    //消息类型 1私聊 2群聊 3广播
	Media      int    //消息类型 1文字 2图片 3音频 4表情包
	Content    string //消息内容
	CreateTime uint64 // 创建时间
	ReadTime   uint64 // 读取时间
	Pic        string
	Url        string
	Desc       string
	Amount     int //其他数字统计
}

func (table *Message) TableName() string {
	return "message"
}

// Node结构体表示一个客户端连接节点
type Node struct {
	Conn          *websocket.Conn //维护了与客户端的 WebSocket 连接。
	Addr          string          //客户端地址
	FirstTime     uint64          //首次连接时间
	HeartbeatTime uint64          //心跳时间
	LoginTime     uint64          //登录时间
	DataQuene     chan []byte     //用于存储需要发送给客户端的消息，是一个消息队列（类型是 chan []byte）。
	GroupSets     set.Interface   //用于存储客户端加入的群组集合，便于群聊时管理群成员。

}

// 映射关系
var ClientMap map[int64]*Node = make(map[int64]*Node, 0) //一个全局映射，用于存储用户ID与Node节点之间的关系，以便根据用户 ID 查找对应的 WebSocket 连接。

// 读写锁
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query() //获取用户id
	Id := query.Get("userId")
	userId, _ := strconv.ParseInt(Id, 10, 64)
	// msgType := query.Get("type")
	// targetId := query.Get("targetId")
	// context := query.Get("context")
	isvalida := true
	// 将HTTP连接升级为websocket连接
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 获取conn
	currentTime := uint64(time.Now().Unix())
	node := &Node{
		Conn:          conn,
		Addr:          conn.RemoteAddr().String(), //客户端地址
		HeartbeatTime: currentTime,                //心跳时间
		LoginTime:     currentTime,                //登陆时间
		DataQuene:     make(chan []byte),
		GroupSets:     set.New(set.ThreadSafe),
	}

	//用户和node绑定并加锁
	rwLocker.RLock()
	ClientMap[userId] = node
	rwLocker.RUnlock()
	// 发送逻辑
	go sendProc(node) //负责发送消息给客户端。
	//接收逻辑
	go receiveProc(node) // 负责接收来自客户端的消息。
	// 加入在线用户到缓存
	SetUserOnlineInfo("online_"+Id, []byte(node.Addr), time.Duration(viper.GetInt("timeout.RedisOnlineTime"))*time.Hour)
	// sendMsg(userId, []byte("欢迎来到聊天系统"))

}

func sendProc(node *Node) {
	for {
		select {
		//该函数在Node的 DataQuene 中不断读取消息，然后通过 WebSocket 发送给客户端。
		case data := <-node.DataQuene:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func receiveProc(node *Node) {
	for {
		//该函数不断从 WebSocket 连接中读取消息，然后调用 broadMsg 将消息广播给其他客户端。
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		msg := Message{}
		err2 := json.Unmarshal(data, &msg)
		if err2 != nil {
			fmt.Println(err2)
		}
		if msg.Type == 3 {
			currentTime := uint64(time.Now().Unix())
			node.Heartbeat(currentTime)
		} else {
			// dispatch(data)
			broadMsg(data) // 将消息广播到局域网
			fmt.Println("[ws] recvproc <<<<<<<", string(data))
		}

	}
}

var udpsendChan chan []byte = make(chan []byte, 1024) // 是一个带有缓冲区的通道，用于存储需要通过 UDP 发送的消息。

func broadMsg(data []byte) {
	udpsendChan <- data
}

func init() {
	go udpSendProc() // UDP发送消息
	go udprecvProc() // UDP接收消息
	fmt.Println("init groutine")
}

func udpSendProc() {
	// 创建UDP连接
	con, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	//从 udpsendChan 读取数据并通过 UDP 发送。
	for {
		select {
		case data := <-udpsendChan:
			fmt.Println("udpsendproc data :>>>>", string(data))
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func udprecvProc() {
	//在本地端口监听 UDP 消息。
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer con.Close()
	//接收消息并调用dispatch方法来进一步处理。
	for {
		var buf [512]byte
		n, err := con.Read(buf[0:]) // 调用Read方法，将接收到的数据存储到buf中。
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("udprecvproc data :", string(buf[0:n]))
		dispatch(buf[0:n])
	}
}

// 后端调度逻辑处理
func dispatch(data []byte) {
	//创建一个空的message对象
	msg := Message{}
	//将接收到的 data 反序列化为 Message 对象。
	msg.CreateTime = uint64(time.Now().Unix())
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println("err的错误:", err)
		return
	}
	switch msg.Type {
	case 1: //私信
		sendMsg(int64(msg.Targetid), data)
	case 2: //群发
		sendGroupMsg(int64(msg.Targetid), data) // 发送的群ID，消息内容
		// case 3: //广播
		// 	sendAllMsg()
	}
}

func sendMsg(userId int64, msg []byte) {
	rwLocker.RLock()
	//使用userId从ClientMap中找到对应的 Node
	node, ok := ClientMap[userId]
	rwLocker.RUnlock()
	// if ok {
	// 	node.DataQuene <- msg
	// }
	jsonMsg := Message{}
	merr := json.Unmarshal(msg, &jsonMsg)
	if merr != nil {
		fmt.Println("转换消息是否有问题》》》》", merr)
	}
	ctx := context.Background()
	targetIdStr := strconv.Itoa(int(userId))
	userIdStr := strconv.Itoa(int(jsonMsg.Userid))
	jsonMsg.CreateTime = uint64(time.Now().Unix())
	r, err := utils.Red.Get(ctx, "online_"+userIdStr).Result()
	if err != nil {
		fmt.Println("redis获取键名是否有误：》》》》》", err)
	}
	if r != "" {
		if ok {
			//将消息放入 DataQuene：找到后将消息放入 DataQuene，供 sendProc 发送给用户
			fmt.Println("sendMsg >>>> userID:", userId, " msg:", string(msg))
			node.DataQuene <- msg
		}
	}
	var key string
	if userId > int64(jsonMsg.Userid) {
		key = "msg" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg" + targetIdStr + "_" + userIdStr
	}
	res, err := utils.Red.ZRevRange(ctx, key, 0, -1).Result()
	if err != nil {
		fmt.Println(err)
	}
	score := float64(cap(res)) + 1
	ress, e := utils.Red.ZAdd(ctx, key, &redis.Z{score, msg}).Result()
	if e != nil {
		fmt.Println(e)
	}
	fmt.Println(ress)
}

// 群发
func sendGroupMsg(userId int64, msg []byte) {
	RtMsg := Message{}
	json.Unmarshal(msg, &RtMsg)
	fmt.Println("开始群发消息")
	userIds := SearchUserByGroupId(uint(userId))
	for _, v := range userIds {
		sendMsg(int64(v), msg)
	}
}

// 获取缓存里面的消息
func RedisMsg(userIdA, userIdB, start, end int, isRev bool) []string {
	rwLocker.RLock()
	rwLocker.RUnlock()

	// 创建一个ctx获取上下文，用于执行redis操作
	/*
		context.Background() 是最初的、最顶层的 Context，它通常用于根级别的上下文，
		比如应用程序的初始化阶段，或者启动新的 goroutine 时，不需要传递任何附加信息的情况。
	*/
	ctx := context.Background()
	// int型转换为string类型
	userIdStr := strconv.Itoa(userIdA)
	targetIdStr := strconv.Itoa(userIdB)
	// 定义一个key变量，用于存储Redis的键
	var key string
	if userIdA > userIdB {
		key = "msg_" + userIdStr + "_" + targetIdStr
	} else {
		key = "msg_" + targetIdStr + "_" + userIdStr
	}
	// redis的一种数据结构：有序集合。ZRange和ZRevRange方法是对其按照分数进行升序或降序排列
	var rels []string
	var err error
	if isRev {
		rels, err = utils.Red.ZRange(ctx, key, int64(start), int64(end)).Result()
	} else {
		rels, err = utils.Red.ZRevRange(ctx, key, int64(start), int64(end)).Result()
	}

	if err != nil {
		fmt.Println(err) // 没有找到记录
	}
	// 发送推送消息
	/**
	// 后台通过websoket 推送消息
	for _, val := range rels {
		fmt.Println("sendMsg >>> userID: ", userIdA, "  msg:", val)
		node.DataQueue <- []byte(val)
	}**/
	return rels

}

// 更新用户心跳
func (node *Node) Heartbeat(currentTime uint64) {
	node.HeartbeatTime = currentTime
	return
}

// 清理超时连接
func CleanConnection(param interface{}) (result bool) {
	result = true
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("cleanConnection err", r)
		}
	}()
	currenTime := uint64(time.Now().Unix())
	for i := range ClientMap {
		node := ClientMap[i]
		if node.IsHeartbeatTimeOut(currenTime) {
			fmt.Println("心跳超时。。。。关闭连接", node)
			node.Conn.Close()
		}
	}
	return result
}

// 用户心跳是否超时
func (node *Node) IsHeartbeatTimeOut(currentTime uint64) (timeout bool) {
	if node.HeartbeatTime+viper.GetUint64("timeout.HeartbeatMaxTime") <= currentTime {
		fmt.Println("心跳超时。。。自动下线", node)
		timeout = true
	}
	return
}

// 设置在线用户到redis缓存
func SetUserOnlineInfo(key string, val []byte, timeTTL time.Duration) {
	ctx := context.Background()
	err := utils.Red.Set(ctx, key, val, timeTTL).Err()
	if err != nil {
		fmt.Println("设置在线用户到redis缓存有误：》》》》", err)
	}
}
