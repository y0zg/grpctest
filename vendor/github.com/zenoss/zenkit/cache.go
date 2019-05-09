package zenkit

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/cache"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func CacheExpiration() time.Duration {
	raw := viper.GetString(GCMemstoreTTLConfig)
	dur, err := time.ParseDuration(raw)
	if err != nil {
		logrus.WithError(err).WithField("cache_duration", raw).Error("Could not parse cache duration")
	}
	return dur
}

func NewCache() *cache.Codec {
	return NewCacheId(0)
}

func NewCacheId(dbid int) *cache.Codec {
	codec := &cache.Codec{
		Marshal:   json.Marshal,
		Unmarshal: json.Unmarshal,
	}
	if rds := NewRedisRingId(dbid); rds != nil {
		codec.Redis = rds
	} else {
		maxLen := viper.GetInt(GCMemstoreLocalMaxLen)
		codec.UseLocalCache(maxLen, CacheExpiration())
	}
	return codec
}

func NewRedisRing() *redis.Ring {
	return NewRedisRingId(0)
}

func NewRedisRingId(dbid int) *redis.Ring {
	redisAddrs := viper.GetStringSlice(GCMemstoreAddressConfig)
	if len(redisAddrs) == 0 {
		return nil
	}
	addrRing := make(map[string]string)
	for i, addr := range redisAddrs {
		addrRing[fmt.Sprintf("%d", i+1)] = addr
	}
	return redis.NewRing(&redis.RingOptions{
		Addrs: addrRing,
		DB: dbid,
	})
}
