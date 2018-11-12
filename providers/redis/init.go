package redis

import (
	"os"
	"strings"

	"code.byted.org/gopkg/env"
	"code.byted.org/gopkg/logs"
	"code.byted.org/learning_fe/go_modules/redis"
)

// InitRedis init Redis client
func InitRedis() (*redis.CacheWrapper, err) {
	var err error

	psm := Config.Redis

	prefix := ""
	if !env.IsProduct() {
		prefix = "learning_open:"
	}

	if strings.HasPrefix(psm, "toutiao.") {
		if Redis, err := redis.NewCacheWrapperWithPrefix(psm, prefix); err != nil {
			logs.Error("cache.NewCacheWrapper error %s", err)
			logs.Stop()
			os.Exit(-1)
			return nil, err
		}
		return Redis, nil
	} else {
		servers := []string{psm}
		if Redis, err = redis.NewCacheWrapperWithServersWithPrefix("toutiao.redis.web_dev", servers, prefix); err != nil {
			logs.Error("cache.NewCacheWrapper error %s", err)
			logs.Stop()
			os.Exit(-1)
			return nil, err
		}
		return Redis, nil
	}
}
