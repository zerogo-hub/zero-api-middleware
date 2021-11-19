package csrf

import (
	"crypto/subtle"
	"net/http"

	zeroapi "github.com/zerogo-hub/zero-api"
	zeroctx "github.com/zerogo-hub/zero-api/context"

	zerocrypto "github.com/zerogo-hub/zero-helper/crypto"
	zerorandom "github.com/zerogo-hub/zero-helper/random"
)

// New 为请求填充 csrf
// 需要启用 cache 功能
func New(opts ...Option) zeroapi.Handler {
	opt := defaultOption()
	if len(opts) > 0 {
		option := opts[0]
		opt.replace(option)
	}

	return func(ctx zeroapi.Context) {
		if opt.IgnoreFunc != nil && opt.IgnoreFunc(ctx) {
			return
		}

		token := zerocrypto.HmacMd5(zerorandom.String(32), opt.Key)
		ctx.SetCookie(
			opt.CookieName,
			token,
			zeroctx.WithCookieMaxAge(opt.CookieMaxAge),
			zeroctx.WithCookieDomain(opt.CookieDomain),
			zeroctx.WithCookiePath(opt.CookiePath),
			zeroctx.WithCookieHTTPOnly(opt.CookieHTTPOnly),
		)
	}
}

// Verify 验证 csrf
func Verify(opts ...Option) zeroapi.Handler {
	opt := defaultOption()
	if len(opts) > 0 {
		option := opts[0]
		opt.replace(option)
	}

	return func(ctx zeroapi.Context) {
		method := ctx.Method()
		for _, requiredMethod := range opt.Methods {
			if method == requiredMethod {
				if !opt.verify(ctx) {
					ctx.Stopped()
					ctx.SetHTTPCode(http.StatusBadRequest)
					ctx.App().Logger().Errorf("invalid csrf token, method: %s, path: %s", ctx.Method(), ctx.Path())
				}
				break
			}
		}
	}
}

func (opt *Option) verify(ctx zeroapi.Context) bool {
	clientToken := opt.clientToken(ctx)
	if len(clientToken) == 0 {
		return false
	}

	cookieToken := opt.cookieToken(ctx)
	if len(cookieToken) == 0 {
		return false
	}

	result := subtle.ConstantTimeCompare([]byte(clientToken), []byte(cookieToken)) == 1
	return result
}

func (opt *Option) clientToken(ctx zeroapi.Context) string {
	// 顺序: query/body/header

	token := ctx.Get(opt.QueryName)
	if len(token) == 0 {
		token = ctx.Query(opt.BodyName)
		if len(token) == 0 {
			token = ctx.Header(opt.HeaderName)
		}
	}

	return token
}

func (opt *Option) cookieToken(ctx zeroapi.Context) string {

	token, err := ctx.Cookie(
		opt.CookieName,
	)

	if err != nil {
		return ""
	}

	return token
}
