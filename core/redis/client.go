package redis

import (
	"context"
	"fmt"
	"time"
	"uwwolf/util"

	"github.com/redis/go-redis/v9"
)

var client *redis.ClusterClient

func ConnectRedis() *redis.ClusterClient {
	if client != nil {
		return client
	}

	fmt.Println("Connecting to redis...")
	client = redis.NewFailoverClusterClient(&redis.FailoverOptions{
		MasterName:      util.Config().Redis.MasterName,
		SentinelAddrs:   util.Config().Redis.SentinelAddresses,
		Username:        util.Config().Redis.Username,
		Password:        util.Config().Redis.Password,
		PoolSize:        10,
		MinRetryBackoff: 1 * time.Second,
		MaxRetryBackoff: 5 * time.Second,
		MaxRetries:      10,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Redis is connected!")

	return client
}
