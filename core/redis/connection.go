package redis

import (
	"context"
	"log"
	"time"
	"uwwolf/config"

	"github.com/redis/go-redis/v9"
)

func Connect(config config.Redis) *redis.ClusterClient {
	client := redis.NewFailoverClusterClient(&redis.FailoverOptions{
		MasterName:      config.MasterName,
		SentinelAddrs:   config.SentinelAddresses,
		Username:        config.Username,
		Password:        config.Password,
		PoolSize:        config.PollSize,
		ConnMaxIdleTime: 0,
		MinRetryBackoff: 5 * time.Second,
		MaxRetryBackoff: 10 * time.Second,
		MaxRetries:      10,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Panic(err)
	}

	return client
}
