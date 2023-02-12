package database

import (
	"github.com/redis/go-redis/v9"
)

var cache *redis.Client

func InitRedis() {
	cache = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func GetRedis() *redis.Client {
	return cache
}
