package load

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

var cache *redis.Client

func InitRedis() {
	r := Conf.Redis
	cache = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", r.Host, r.Port),
		Password: r.Password, // no password set
		DB:       r.Database, // use default DB
	})
}

func GetRedis() *redis.Client {
	return cache
}
