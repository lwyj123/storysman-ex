package redis

import (
	"sync"

	"code.byted.org/learning_fe/go_modules/redis"
)

// ArticleClient article的thrift Client包装
type Redis struct {
	redis *redis.CacheWrapper
}

var _redis *Redis
var _redisOnce sync.Once

// ArticleClientInstance 返回 ArticleClient 单例对象
func RedisInstance() *Redis {
	_redisOnce.Do(func() {
		redis, _ := InitRedis()
		_redis = &Redis{
			redis: redis,
			// cache??
		}
	})
	return _redis
}

// GetClient 获取redis的client
func (r *Redis) GetClient() (*goredis.Client, error) {
	return r.redis.GetClient(), nil
}
