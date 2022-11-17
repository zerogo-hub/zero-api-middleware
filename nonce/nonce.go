package nonce

import (
	"net/http"

	zeroapi "github.com/zerogo-hub/zero-api"
	zerocache "github.com/zerogo-hub/zero-helper/cache"
)

// New 判断 nonce 是否有效，在一段时间内其不可重复
// 需要启用 cache 功能
func New(cache zerocache.Cache, opts ...Option) zeroapi.Handler {
	opt := defaultOption()
	if len(opts) > 0 {
		opt = opts[0]
	}

	return func(ctx zeroapi.Context) {
		if !opt.Enable {
			return
		}

		s := ctx.Query(opt.Field)
		key := opt.PrefixNonce + s
		exist, _ := cache.Exists(key)
		if exist {
			ctx.Stopped()
			ctx.SetHTTPCode(http.StatusBadRequest)
			ctx.App().Logger().Errorf("repeated nonce, method: %s, path: %s, ip: %s", ctx.Method(), ctx.Path(), ctx.IP())
			return
		}

		if err := cache.SetEx(key, "1", opt.Expire); err != nil {
			ctx.App().Logger().Errorf("cache nonce failed, err: %s", err.Error())
		}
	}
}
