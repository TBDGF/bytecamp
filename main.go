package main

import (
	"bytedance/auth"
	"bytedance/member"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func main() {
	r := gin.Default()
	// 连接数据库
	database, err := sqlx.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/camp")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}
	fmt.Println("open mysql success.")

	g := r.Group("/api/v1")
	// 登录模块
	{
		auth.Login(g, database)
		auth.Logout(g)
		auth.Whoami(g, database)
	}

	// 成员模块
	{
		member.GetMember(g, database)
		member.Create(g, database)
	}

	// 排课模块
	{

	}

	// 抢课模块
	{

	}

	err = r.Run(":80")
	if err != nil {
		return
	}
}
