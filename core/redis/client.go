package redis

import (
	"context"
	"fmt"
	"time"
	"uwwolf/util"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client

func ConnectRedis() *redis.Client {
	if client != nil {
		return client
	}

	fmt.Println("Connecting to redis...")
	client = redis.NewClient(&redis.Options{
		Addr:            util.Config().Redis.Hosts[0],
		Password:        util.Config().Redis.Password,
		MinRetryBackoff: 1 * time.Second,
		MaxRetryBackoff: 5 * time.Second,
		MaxRetries:      10,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	fmt.Println("Connected!")
	return client
}
