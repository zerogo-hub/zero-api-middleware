package limiter

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"

	zeroapi "github.com/zerogo-hub/zero-api"
)

var (
	limiter *rate.Limiter
)

// New 全局限流器，每隔 every 时间放入一个令牌，满 burst 个令牌后不放入新令牌
func New(every time.Duration, burst int) zeroapi.Handler {

	// 每隔 every 时间放入 1 个，初始放入 burst 个
	limiter = rate.NewLimiter(rate.Every(every), burst)

	// 每秒 10 个，上限 100
	// limiter = rate.NewLimiter(rate.Limit(10), 100)

	return func(ctx zeroapi.Context) {
		if ctx.Method() == http.MethodOptions {
			return
		}

		if !limiter.Allow() {
			ctx.Stopped()
			ctx.SetHTTPCode(http.StatusForbidden)
			ctx.App().Logger().Errorf("global limiter, method: %s, path: %s, ip: %s", ctx.Method(), ctx.Path(), ctx.IP())
			return
		}
	}
}

// SetLimit 动态修改放入令牌的速率
func SetLimit(every time.Duration) {
	if limiter != nil {
		limiter.SetLimit(rate.Every(every))
	}
}

// SetBurst 动态修改令牌桶大小
func SetBurst(burst int) {
	if limiter != nil {
		limiter.SetBurst(burst)
	}
}
