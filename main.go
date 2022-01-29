package main

import (
	"bytedance/config"
	"bytedance/router"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	r := gin.Default()
	config.InitDB()
	router.RegisterRouter(r)
	r.Run(":80")
}
