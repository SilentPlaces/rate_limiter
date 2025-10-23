package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/SilentPlaces/rate_limiter/internal/application/ports"
	"github.com/redis/go-redis/v9"
)

type RedisAdapter struct {
	client           *redis.Client
	logger           ports.Logger
	operationTimeout time.Duration
}

func NewRedisAdapter(client *redis.Client, logger ports.Logger, operationTimeoutSeconds int) ports.LimiterScore {
	return &RedisAdapter{
		client:           client,
		logger:           logger,
		operationTimeout: time.Duration(operationTimeoutSeconds) * time.Second,
	}
}

func (r *RedisAdapter) Get(ctx context.Context, key string) interface{} {
	ctx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()
	cmd := r.client.Get(ctx, key)
	val, err := cmd.Result()
	if err != nil {
		r.logger.Error(fmt.Sprintf("RedisAdapter:Get error for key: %s", key), ports.Field{Key: "error", Val: err})
		return err
	}
	return val
}

func (r *RedisAdapter) Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error {
	ctx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()
	cmd := r.client.Set(ctx, key, value, time.Duration(ttlSeconds)*time.Second)
	if err := cmd.Err(); err != nil {
		r.logger.Error(fmt.Sprintf("RedisAdapter:Set error for key: %s", key), ports.Field{Key: "error", Val: err})
		return err
	}
	return nil
}

func (r *RedisAdapter) Incr(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()
	cmd := r.client.Incr(ctx, key)
	if err := cmd.Err(); err != nil {
		r.logger.Error(fmt.Sprintf("RedisAdapter:Incr error for key: %s", key), ports.Field{Key: "error", Val: err})
		return err
	}
	return nil
}

func (r *RedisAdapter) Eval(ctx context.Context, script string, keys []string, args ...[]interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, r.operationTimeout)
	defer cancel()
	flatArgs := make([]interface{}, 0)
	for _, arr := range args {
		flatArgs = append(flatArgs, arr...)
	}
	cmd := r.client.Eval(ctx, script, keys, flatArgs...)
	res, err := cmd.Result()
	if err != nil {
		r.logger.Error("RedisAdapter:Eval error", ports.Field{Key: "error", Val: err})
		return nil, err
	}
	return res, nil
}
