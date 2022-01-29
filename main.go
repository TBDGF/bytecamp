package main

import (
	"bytedance/config"
	"bytedance/router"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()
	config.InitDB()          // 数据库
	router.RegisterRouter(r) // 接口
	r.Run(":80")
}
