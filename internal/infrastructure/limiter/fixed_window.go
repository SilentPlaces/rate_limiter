package limiter

import (
	"context"
	"fmt"

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

func (f *FixedWindowLimiter) Allow(ctx context.Context, key string, cfg config.AlgorithmConfig) (bool, error) {
	fixedCfg, ok := cfg.(config.FixedWindowConfig)
	if !ok {
		return false, fmt.Errorf("invalid config type for FixedWindowLimiter, got %T", cfg)
	}

	res, err := f.score.Eval(ctx, f.luaScript, []string{key}, []interface{}{fixedCfg.Window, fixedCfg.Limit})
	if err != nil {
		return false, err
	}

	allowed, ok := res.(int64)
	return ok && allowed == 1, nil
}
