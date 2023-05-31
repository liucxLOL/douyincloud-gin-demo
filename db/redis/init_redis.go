package redis

import (
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisInstance *redis.Client

func InitRedis() {
	startingTime := time.Now().UTC()
	user := os.Getenv("REDIS_USERNAME")
	pwd := os.Getenv("REDIS_PASSWORD")
	addr := os.Getenv("REDIS_ADDRESS")

	redisInstance = redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: user,
		Password: pwd, // no password set
		DB:       0,   // use default DB
	})
	endingTime := time.Now().UTC()
	fmt.Println("finish init redis, duration is: ", endingTime.Sub(startingTime))
}

func GetRedis() *redis.Client {
	return redisInstance
}
