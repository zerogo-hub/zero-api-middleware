// Package throttle 限流
package throttle

import (
	"math"
	"net/http"
	"strconv"

	throttled "github.com/throttled/throttled/v2"
	memstore "github.com/throttled/throttled/v2/store/memstore"
	zeroapi "github.com/zerogo-hub/zero-api"
)

// VaryBy 为请求生成唯一值
type VaryBy interface {
	Key(*http.Request) string
}

// Config 配置
type Config struct {
	// LimitHandler 发生限制时调用
	LimitHandler zeroapi.Handler

	// ErrHandler 发生错误时调用
	ErrHandler func(ctx zeroapi.Context, err error)

	// PerMin 每个请求每分钟请求的个数，> 0
	PerMin int

	// Burst 每一个请求每分钟允许额外请求的个数，> 0
	Burst int

	// VaryBy 为请求生成唯一值
	VaryBy VaryBy
}

func New(c *Config) zeroapi.Handler {

	if c == nil {
		return nil
	}

	if c.LimitHandler == nil {
		c.LimitHandler = DefaultLimitHandler
	}

	if c.ErrHandler == nil {
		c.ErrHandler = DefaultErrHandler
	}

	if c.PerMin == 0 {
		c.PerMin = 20
	}

	if c.Burst == 0 {
		c.Burst = 1
	}

	if c.VaryBy == nil {
		c.VaryBy = DefaultVaryBy()
	}

	store, err := memstore.New(65536)
	if err != nil {
		return nil
	}

	quota := throttled.RateQuota{
		MaxRate:  throttled.PerMin(c.PerMin),
		MaxBurst: c.Burst,
	}

	rateLimiter, err := throttled.NewGCRARateLimiter(store, quota)
	if err != nil {
		return nil
	}

	return func(ctx zeroapi.Context) {
		var key string
		if c.VaryBy != nil {
			key = c.VaryBy.Key(ctx.Request())
		}

		limited, context, err := rateLimiter.RateLimit(key, 1)
		if err != nil {
			c.ErrHandler(ctx, err)
			return
		}

		setRateLimitHeaders(ctx, context)

		if limited {
			c.LimitHandler(ctx)
			return
		}
	}
}

// DefaultLimitHandler 默认的发生限制时调用
func DefaultLimitHandler(ctx zeroapi.Context) {
	ctx.SetHTTPCode(http.StatusTooManyRequests)
	ctx.Message(http.StatusTooManyRequests, "too many requests")
	ctx.Stopped()
}

// DefaultErrHandler 默认的发生错误时调用
func DefaultErrHandler(ctx zeroapi.Context, err error) {
	ctx.SetHTTPCode(http.StatusInternalServerError)
	ctx.Message(http.StatusInternalServerError, "internal server error")
	ctx.Stopped()
}

// DefaultVaryBy 默认的 key 生成器
func DefaultVaryBy() VaryBy {
	return &throttled.VaryBy{
		Separator:  "_",
		RemoteAddr: true,
		Method:     true,
		Path:       true,
	}
}

func setRateLimitHeaders(ctx zeroapi.Context, context throttled.RateLimitResult) {
	if v := context.Limit; v >= 0 {
		// 最大访问次数
		ctx.SetHeader("X-RateLimit-Limit", strconv.Itoa(v))
	}

	if v := context.Remaining; v >= 0 {
		// 剩余访问次数
		ctx.SetHeader("X-RateLimit-Remaining", strconv.Itoa(v))
	}

	if v := context.ResetAfter; v >= 0 {
		vi := int(math.Ceil(v.Seconds()))
		// 时间过后，访问次数重置
		ctx.SetHeader("X-RateLimit-Reset", strconv.Itoa(vi))
	}

	if v := context.RetryAfter; v >= 0 {
		vi := int(math.Ceil(v.Seconds()))
		// 需要等待多长时间之后才能继续发送请求
		ctx.SetHeader("Retry-After", strconv.Itoa(vi))
	}
}
