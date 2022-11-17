package limiter

import (
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	zeroapi "github.com/zerogo-hub/zero-api"
)

// NewIP 针对 ip 的限流器，每隔 every 时间放入一个令牌，满 burst 个令牌后不放入新令牌
func NewIP(every time.Duration, burst int) zeroapi.Handler {

	i := newIPRateLimiter(rate.Every(every), burst)

	return func(ctx zeroapi.Context) {
		if ctx.Method() == http.MethodOptions {
			return
		}

		ipStr := ctx.IP()
		l := i.getLimiter(ipStr)

		if !l.Allow() {
			ctx.Stopped()
			ctx.SetHTTPCode(http.StatusForbidden)
			ctx.App().Logger().Errorf("ip limiter: %s, method: %s, path: %s, ip: %s", ipStr, ctx.Method(), ctx.Path(), ctx.IP())
			return
		}
	}
}

type ipRateLimter struct {
	ips   map[string]*rate.Limiter
	lock  *sync.RWMutex
	r     rate.Limit
	burst int
}

func newIPRateLimiter(r rate.Limit, burst int) *ipRateLimter {
	limiter := &ipRateLimter{
		ips:   make(map[string]*rate.Limiter),
		lock:  &sync.RWMutex{},
		r:     r,
		burst: burst,
	}

	return limiter
}

func (limiter *ipRateLimter) addIP(ipStr string) *rate.Limiter {
	limiter.lock.Lock()
	defer limiter.lock.Unlock()

	l := rate.NewLimiter(limiter.r, limiter.burst)
	limiter.ips[ipStr] = l

	return l
}

func (limiter *ipRateLimter) getLimiter(ipStr string) *rate.Limiter {
	limiter.lock.RLock()

	l, exists := limiter.ips[ipStr]

	if !exists {
		limiter.lock.RUnlock()
		return limiter.addIP(ipStr)
	}

	limiter.lock.RUnlock()

	return l
}
