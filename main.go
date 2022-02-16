package main

import (
	"bytedance/db"
	"bytedance/redis_server"
	"bytedance/router"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func main() {
	r := gin.Default()
	router.RegisterRouter(r) // 路由
	db.InitDB()              // 数据库
	redis_server.InitRedis() // redis连接
	// 日志输出到文件
	file := "./" + "message" + ".log"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	r.Run(":80")
}
