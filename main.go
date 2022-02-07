package main

import (
	"bytedance/db"
	"bytedance/router"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()
	db.InitDB()              // 数据库
	router.RegisterRouter(r) // 路由
	r.Run(":9090")
}
