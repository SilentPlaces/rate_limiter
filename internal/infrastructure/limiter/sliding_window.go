package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/SilentPlaces/rate_limiter/internal/application/ports"
	"github.com/SilentPlaces/rate_limiter/internal/domain/config"
	"github.com/SilentPlaces/rate_limiter/internal/domain/errors"
)

type SlidingWindowLimiter struct {
	score     ports.LimiterScore
	luaScript string
}

func NewSlidingWindowLimiter(score ports.LimiterScore, luaScript string) ports.RateLimiter {
	return &SlidingWindowLimiter{luaScript: luaScript, score: score}
}

func SlidingWindowLimiterFactory(score ports.LimiterScore, luaScript string) ports.RateLimiter {
	return NewSlidingWindowLimiter(score, luaScript)
}

func (s *SlidingWindowLimiter) Allow(ctx context.Context, key string, cfg config.AlgorithmConfig) (bool, error) {
	slidingConfig, ok := cfg.(config.SlidingWindowConfig)
	if !ok {
		return false, errors.NewRateLimiterError(errors.ErrInvalidConfig.Code,
			"invalid config type for SlidingWindowLimiter",
			fmt.Errorf("invalid config type for SlidingWindowLimiter, got %T", cfg))
	}

	now := time.Now().Unix()

	res, err := s.score.Eval(ctx,
		s.luaScript,
		[]string{key}, []interface{}{slidingConfig.Window, slidingConfig.Limit, now},
	)
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
