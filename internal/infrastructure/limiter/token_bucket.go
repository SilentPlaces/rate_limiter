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

func (t *TokenBucketLimiter) Allow(ctx context.Context, key string, cfg config.AlgorithmConfig) (bool, error) {
	tokenCfg, ok := cfg.(config.TokenBucketConfig)
	if !ok {
		return false, fmt.Errorf("invalid config type for TokenBucketLimiter, got %T", cfg)
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
		return false, err
	}

	result, ok := res.([]interface{})
	if !ok || len(result) < 1 {
		return false, fmt.Errorf("unexpected lua script response")
	}

	allowed, ok := result[0].(int64)
	return ok && allowed == 1, nil
}
