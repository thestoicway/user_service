package sessiondatabase

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type KeyValueDatabase interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error)
	Get(ctx context.Context, key string) (value interface{}, err error)
	Delete(ctx context.Context, key string) (err error)
}

type redisDatabase struct {
	redis  *redis.Client
	logger *zap.SugaredLogger
}

func NewRedisDatabase(redis *redis.Client, logger *zap.SugaredLogger) KeyValueDatabase {
	return &redisDatabase{
		redis:  redis,
		logger: logger,
	}
}

// Delete implements KeyValueDatabase.
func (r *redisDatabase) Delete(ctx context.Context, key string) (err error) {
	result := r.redis.Del(ctx, key)

	return result.Err()
}

// Get implements KeyValueDatabase.
func (r *redisDatabase) Get(ctx context.Context, key string) (value interface{}, err error) {
	result := r.redis.Get(ctx, key)

	err = result.Err()

	if err != nil {
		return nil, err
	}

	return result.Val(), nil
}

// Set implements KeyValueDatabase.
func (r *redisDatabase) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error) {
	result := r.redis.Set(ctx, key, value, expiration)

	return result.Err()
}
