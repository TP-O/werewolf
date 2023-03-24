package redis

import (
	"context"
	"time"
	"uwwolf/config"

	"github.com/redis/go-redis/v9"
)

func Connect() *redis.ClusterClient {
	client := redis.NewFailoverClusterClient(&redis.FailoverOptions{
		MasterName:      config.Redis().MasterName,
		SentinelAddrs:   config.Redis().SentinelAddresses,
		Username:        config.Redis().Username,
		Password:        config.Redis().Password,
		PoolSize:        config.Redis().PollSize,
		ConnMaxIdleTime: 0,
		MinRetryBackoff: 5 * time.Second,
		MaxRetryBackoff: 10 * time.Second,
		MaxRetries:      10,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}

	return client
}
