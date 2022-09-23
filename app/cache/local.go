package cache

import (
	"time"

	"github.com/jellydator/ttlcache/v3"

	"uwwolf/app/types"
	"uwwolf/config"
)

var local *ttlcache.Cache[types.CacheKey, any]

func init() {
	ttl, _ := time.ParseDuration(config.Cache.LocalTTL)

	local = ttlcache.New(
		ttlcache.WithTTL[string, any](ttl),
	)

	go local.Start()
}

func Local() *ttlcache.Cache[types.CacheKey, any] {
	return local
}
