// Package csrf 跨站请求伪造
package csrf

import (
	"crypto/subtle"
	"errors"
	"net/http"

	zeroapi "github.com/zerogo-hub/zero-api"
	"github.com/zerogo-hub/zero-helper/crypto"
	"github.com/zerogo-hub/zero-helper/random"
)

// Config 跨站请求伪造 配置
type Config struct {
	// Key 对 cookie 进行签名验证
	Key string

	// CookieName token 在 cookie 中的名字
	CookieName string
	// HeaderName token 在 header 中的名字
	HeaderName string
	// BodyName token 在 body 中的名字
	BodyName string
	// QueryName token 在 query 中的名字
	QueryName string

	// CookieMaxAge cookie 有效时间，秒，默认为 24小时
	CookieMaxAge int
	// CookieDomain ..
	CookieDomain string
	// CookiePath ..
	CookiePath string
	// CookieHTTPOnly
	CookieHTTPOnly bool

	// IgnoreFunc 忽略检测，返回 true 表示不检测
	IgnoreFunc func(ctx zeroapi.Context) bool
}

var defaultConfig = &Config{
	Key:            "csrfKey",
	CookieName:     "csrfToken",
	HeaderName:     "x-csrf-token",
	BodyName:       "_csrf",
	QueryName:      "_csrf",
	CookieMaxAge:   24 * 3600,
	CookieHTTPOnly: true,
}

// New 跨站请求伪造
func New(config *Config) zeroapi.Handler {
	if config == nil {
		config = defaultConfig
	}

	c := defaultConfig
	c.init(config)

	return func(ctx zeroapi.Context) {
		// 忽略，无需检查
		if c.IgnoreFunc != nil && c.IgnoreFunc(ctx) {
			return
		}

		token, err := c.tokenFromCookie(ctx)

		// 从 cookie 获取的 token 值验证失败
		if err != nil && err != http.ErrNoCookie {
			onFailed(ctx)
			return
		}

		switch ctx.Method() {
		case zeroapi.MethodGet, zeroapi.MethodHead, zeroapi.MethodOptions, zeroapi.MethodTrace:
			if token == "" {
				token = crypto.HmacMd5(random.String(32), c.Key)
				c.saveTokenToCookie(token, ctx)
			}
		default:
			if !c.check(token, ctx) {
				onFailed(ctx)
				return
			}
		}
	}
}

func onFailed(ctx zeroapi.Context) {
	ctx.SetHTTPCode(http.StatusForbidden)
	ctx.Message(http.StatusForbidden, "invalid csrf token")
	ctx.Stopped()
}

func (c *Config) init(config *Config) {
	if len(config.Key) != 0 {
		c.Key = config.Key
	}
	if len(config.CookieName) != 0 {
		c.CookieName = config.CookieName
	}
	if len(config.HeaderName) != 0 {
		c.HeaderName = config.HeaderName
	}
	if len(config.BodyName) != 0 {
		c.BodyName = config.BodyName
	}
	if len(config.QueryName) != 0 {
		c.QueryName = config.QueryName
	}
	if config.CookieMaxAge >= 0 {
		c.CookieMaxAge = config.CookieMaxAge
	}
	if len(config.CookieDomain) != 0 {
		c.CookieDomain = config.CookieDomain
	}
	if len(config.CookiePath) != 0 {
		c.CookiePath = config.CookiePath
	}

	c.CookieHTTPOnly = config.CookieHTTPOnly
	c.IgnoreFunc = config.IgnoreFunc
}

// tokenFromCookie 从请求中获取 cookie 值
func (c *Config) tokenFromCookie(ctx zeroapi.Context) (string, error) {
	// 获取 cookie，并验证签名
	cookieValue, err := ctx.Cookie(c.CookieName, zeroapi.WithCookieVerify(c.Key))
	if err != nil {
		return "", err
	}

	if cookieValue == "" {
		return "", errors.New("empty cookie")
	}

	return cookieValue, nil
}

// tokenFromRequest 获取请求中的 token 值
func (c *Config) tokenFromRequest(ctx zeroapi.Context) string {
	// query/body/header

	// query
	token := ctx.Get(c.QueryName)
	if token == "" {
		// body
		token = ctx.Query(c.BodyName)
		if token == "" {
			// header
			token = ctx.Header(c.HeaderName)
		}
	}

	return token
}

// check 检查 cookie 中的 token 值与请求中的 token 值是否一致
func (c *Config) check(token string, ctx zeroapi.Context) bool {
	if token == "" {
		return false
	}

	tokenFromRequest := c.tokenFromRequest(ctx)
	if tokenFromRequest == "" {
		return false
	}

	result := subtle.ConstantTimeCompare([]byte(token), []byte(tokenFromRequest)) == 1
	if !result {
		ctx.App().Logger().Errorf("invalid csrf token, token: %s, tokenFromRequest: %s", token, tokenFromRequest)
	}

	return result
}

// saveTokenToCookie 将 token 保存到 cookie 中
func (c *Config) saveTokenToCookie(cookie string, ctx zeroapi.Context) {
	ctx.AddHeader("Vary", "Cookie")

	opts := []zeroapi.CookieOption{
		zeroapi.WithCookieMaxAge(c.CookieMaxAge),
		zeroapi.WithCookiePath(c.CookiePath),
		zeroapi.WithCookieDomain(c.CookieDomain),
		zeroapi.WithCookieHTTPOnly(c.CookieHTTPOnly),
		zeroapi.WithCookieSign(c.Key),
	}

	ctx.SetCookie(
		c.CookieName,
		cookie,
		opts...,
	)
}
