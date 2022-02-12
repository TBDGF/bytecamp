package redis_server

import (
	"github.com/go-redis/redis"
	"log"
)

var client *redis.Client

func InitRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	log.Println(pong, err)
}

func NewClient() *redis.Client {
	return client
}
