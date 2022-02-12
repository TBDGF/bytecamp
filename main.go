package main

import (
	"bytedance/db"
	"bytedance/redis_server"
	"bytedance/router"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()
	db.InitDB()              // 数据库
	redis_server.InitRedis() // redis连接
	router.RegisterRouter(r) // 路由
	r.Run(":80")
}
