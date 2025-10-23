package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/SilentPlaces/rate_limiter/internal/application/ports"
	"github.com/SilentPlaces/rate_limiter/internal/domain/config"
)

type FixedWindowLimiter struct {
	score     ports.LimiterScore
	luaScript string
}

func NewFixedWindowLimiter(score ports.LimiterScore, luaScript string) ports.RateLimiter {
	return &FixedWindowLimiter{
		score:     score,
		luaScript: luaScript,
	}
}

func FixedWindowLimiterFactory(score ports.LimiterScore, luaScript string) ports.RateLimiter {
	return NewFixedWindowLimiter(score, luaScript)
}

func (f *FixedWindowLimiter) Allow(ctx context.Context, key string, cfg config.AlgorithmConfig) (ports.RateLimitInfo, error) {
	fixedCfg, ok := cfg.(config.FixedWindowConfig)
	if !ok {
		return ports.RateLimitInfo{}, fmt.Errorf("invalid config type for FixedWindowLimiter, got %T", cfg)
	}

	res, err := f.score.Eval(ctx, f.luaScript, []string{key}, []interface{}{fixedCfg.Window, fixedCfg.Limit})
	if err != nil {
		return ports.RateLimitInfo{}, err
	}

	result, ok := res.([]interface{})
	if !ok || len(result) < 4 {
		return ports.RateLimitInfo{}, fmt.Errorf("unexpected lua script response")
	}

	allowed, _ := result[0].(int64)
	remaining, _ := result[2].(int64)
	ttl, _ := result[3].(int64)

	var resetTime int64
	if ttl > 0 {
		resetTime = time.Now().Unix() + ttl
	}

	return ports.RateLimitInfo{
		Allowed:   allowed == 1,
		Limit:     fixedCfg.Limit,
		Remaining: int(remaining),
		ResetTime: resetTime,
	}, nil
}
