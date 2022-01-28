package main

import (
	"bytedance/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	g := gin.Default()
	// 连接数据库
	database, err := sqlx.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/camp")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	fmt.Println("open mysql success.")

	// 登录功能
	{
		auth.Login(g, database)
		auth.Logout(g, database)
		auth.Whoami(g, database)
	}

	g.Run(":80")
}
