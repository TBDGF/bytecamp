package redis_server

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
)

var client *redis.Client

func InitRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     "192.168.80.130:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}

func Redis() *redis.Client {
	log.Println("redis_server:", client)
	return client
}
