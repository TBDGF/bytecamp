package redis_server

import (
	"log"
	"testing"
)

func TestRedis(t *testing.T) {
	InitRedis()
	client := Redis()

	err := client.Set("key", "value", 0).Err()
	if err != nil {
		log.Println("err: ", err)
	}
	val, err := client.Get("2").Result()
	if err != nil {
		log.Println("err: ", err)
	}
	log.Printf("%T", val)
}
