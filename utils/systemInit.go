package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB  *gorm.DB
	Red *redis.Client
)

func InitConfig() {
	viper.SetConfigName("a")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("config app inited")
}

func InitMysql() {
	//自定义日志模板，打印SQL语句
	newloger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, //慢SQL阈值,表示执行时间超过 1 秒的 SQL 查询会被记录下来
			LogLevel:      logger.Info, //级别
			Colorful:      true,        //彩色
		},
	)
	dsn := viper.GetString("mysql.dns")
	DB, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newloger})
	fmt.Println("config mysql inited")
	// user := &models.UserBasic{}
	// DB.Find(user)
	// fmt.Println(user)
}

func InitRedis() {

	Red = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.minIdleConns"),
	})
}

const (
	PublishKey = "websocket"
)

// Subscribe和Publish用于‘服务器内部’的消息分发，使得消息可以广播给多个websocket客户端
// Publish 发送消息到Redis
func Publish(ctx context.Context, channel string, msg string) error {
	var err error = Red.Publish(ctx, channel, msg).Err()
	if err != nil {
		fmt.Println(err)
	}
	return err
}

// Subscribe 订阅Redis消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := Red.Subscribe(ctx, channel)
	msg, err := sub.ReceiveMessage(ctx)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return msg.Payload, err
}
