package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/SilentPlaces/rate_limiter/internal/application/ports"
	"github.com/SilentPlaces/rate_limiter/internal/domain/config"
)

type TokenBucketLimiter struct {
	score     ports.LimiterScore
	luaScript string
}

func NewTokenBucketLimiter(score ports.LimiterScore, luaScript string) ports.RateLimiter {
	return &TokenBucketLimiter{
		score:     score,
		luaScript: luaScript,
	}
}

func TokenBucketLimiterFactory(score ports.LimiterScore, luaScript string) ports.RateLimiter {
	return NewTokenBucketLimiter(score, luaScript)
}

func (t *TokenBucketLimiter) Allow(ctx context.Context, key string, cfg config.AlgorithmConfig) (ports.RateLimitInfo, error) {
	tokenCfg, ok := cfg.(config.TokenBucketConfig)
	if !ok {
		return ports.RateLimitInfo{}, fmt.Errorf("invalid config type for TokenBucketLimiter, got %T", cfg)
	}

	now := time.Now().Unix()
	tokensToConsume := 1

	res, err := t.score.Eval(ctx, t.luaScript, []string{key}, []interface{}{
		tokenCfg.Capacity,
		tokenCfg.RefillRate,
		tokensToConsume,
		now,
		tokenCfg.BucketTTL,
	})
	if err != nil {
		return ports.RateLimitInfo{}, err
	}

	result, ok := res.([]interface{})
	if !ok || len(result) < 4 {
		return ports.RateLimitInfo{}, fmt.Errorf("unexpected lua script response")
	}

	allowed, _ := result[0].(int64)
	remaining, _ := result[2].(int64)
	resetTime, _ := result[3].(int64)

	return ports.RateLimitInfo{
		Allowed:   allowed == 1,
		Limit:     tokenCfg.Capacity,
		Remaining: int(remaining),
		ResetTime: resetTime,
	}, nil
}
