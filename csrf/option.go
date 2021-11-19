package csrf

import (
	"net/http"

	zeroapi "github.com/zerogo-hub/zero-api"
)

// Option ..
type Option struct {
	// Key 生成 token 使用的密钥
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

	// Methods 这些方法参与验证
	Methods []string

	// IgnoreFunc 忽略检测，返回 true 表示不检测
	IgnoreFunc func(ctx zeroapi.Context) bool
}

func defaultOption() *Option {
	return &Option{
		Key:            "csrf key",
		CookieName:     "csrfToken",
		HeaderName:     "x-csrf-token",
		BodyName:       "_csrf",
		QueryName:      "_csrf",
		CookieMaxAge:   24 * 3600,
		CookieHTTPOnly: true,
		Methods:        []string{http.MethodPost, http.MethodPut, http.MethodDelete},
	}
}

func (opt *Option) replace(option Option) {
	if len(option.Key) > 0 {
		opt.Key = option.Key
	}
	if len(option.CookieName) > 0 {
		opt.CookieName = option.CookieName
	}
	if len(option.HeaderName) > 0 {
		opt.HeaderName = option.HeaderName
	}
	if len(option.BodyName) > 0 {
		opt.BodyName = option.BodyName
	}
	if len(option.QueryName) > 0 {
		opt.QueryName = option.QueryName
	}
	if option.CookieMaxAge > 0 {
		opt.CookieMaxAge = option.CookieMaxAge
	}
	if option.IgnoreFunc != nil {
		opt.IgnoreFunc = option.IgnoreFunc
	}
	opt.CookieHTTPOnly = option.CookieHTTPOnly
	if option.Methods != nil {
		opt.Methods = option.Methods
	}
}
