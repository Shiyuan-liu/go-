package service

import (
	"fmt"
	"ginchat/utils"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func Upload(c *gin.Context) {
	srcfile, filehead, err := c.Request.FormFile("file")
	if err != nil {
		utils.RespFail(c.Writer, err.Error())
	}
	suffix := ".png"
	ofile := filehead.Filename
	tem := strings.Split(ofile, ".")
	if len(tem) > 1 {
		suffix = "." + tem[len(tem)-1]
	}
	fileName := fmt.Sprintf("%d%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	// os 包：用于与操作系统交互，例如文件操作、环境变量、进程管理等。
	dstFile, err := os.Create("./asset/upload/" + fileName)
	if err != nil {
		utils.RespFail(c.Writer, err.Error())
	}
	// io 包：提供了与数据流和接口相关的抽象，例如读取、写入、复制、缓冲等功能。
	_, err1 := io.Copy(dstFile, srcfile) // io.Copy(x,y)将y中的数据复制到x中
	if err1 != nil {
		utils.RespFail(c.Writer, err1.Error())
	}
	url := "./asset/upload/" + fileName
	utils.RespOk(c.Writer, "发送图片成功", url)

}
