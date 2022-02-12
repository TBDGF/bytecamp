package redis_server

import (
	"bytedance/db"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"testing"
)

func TestRedis(t *testing.T) {
	db.InitDB()
	InitRedis()

	log.Println(GetMemberByID("1"))
}
