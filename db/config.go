package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

var Db *sqlx.DB

// 数据库连接参数
const (
	SQLUser     string = "root"
	SQLPassword string = "123456"
	Host        string = "127.0.0.1"
	port        string = "3306"
	database    string = "camp"
)

func InitDB() {
	var err error
	Db, err = sqlx.Connect("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", SQLUser, SQLPassword, Host, port, database))
	if err != nil {
		log.Println("Init database failed:", err)
		return
	}
	log.Println("Successfully init database.")
}

func NewDB() *sqlx.DB {
	return Db
}
