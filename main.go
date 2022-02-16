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
	//路由
	router.RegisterRouter(r)
	// 数据库
	db.InitDB()
	// redis连接
	redis_server.InitRedis()
	r.Run(":80")
}
