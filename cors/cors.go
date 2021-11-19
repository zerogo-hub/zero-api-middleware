// Package cors 跨域控制
package cors

import (
	"net/http"
	"regexp"
	"strings"

	zeroapi "github.com/zerogo-hub/zero-api"
)

// 简单请求，浏览器直接发出 CORS 请求
//
// 非简单请求，浏览器会增加一次 HTTP 查询请求(预检)，询问服务器，当前网页所在的域名是否在服务器的许可名单之内
// 询问使用的方法是 OPTIONS，带有以下参数:
// Origin: 表示请求来自哪个源
// Access-Control-Request-Method: 该字段必须
// Access-Control-Request-Headers: 使用逗号分隔，指定浏览器发出 CORS 请求会额外发送的头信息字段
//
// 简单请求，需要同时满足以下两个条件
// 1 请求必须是以下三种方法之一: HEAD, GET, POST
// 2 HTTP 头信息不超出以下几种:
// 2-1: Accept
// 2-2: Accept-Language
// 2-3: Content-Language
// 2-4: Last-Event-ID
// 2-5: Content-Type 且仅限于 app/x-www-form-urlencoded、multipart/form-data、text/plain

// Config 跨域配置
type Config struct {
	// AccessControlAllowOrigin (简单请求的响应)请求中 Origin 的值，允许跨域的来源，默认为 "*" ，表示任何来源均可
	AccessControlAllowOrigin []string

	// AccessControlAllowCredentials (简单请求的响应)是否允许浏览器在 CORS 请求中发送 Cookie，默认 false
	// 一般 Cookie 不会包括在 CORS 请求中
	// 传给浏览器时，如果为 false 则不传递该值。如果传递，只能为 true
	AccessControlAllowCredentials bool

	// AccessControlExposeHeaders (简单请求的响应)暴露额外的 header 字段，这样浏览器可以访问更多的 HEADER
	// 默认只能获取 Cache-Control、Content-Language、Content-Length、Content-Type、Expires、Last-Modified、Pragma 这七个字段
	// Content-Type: application/x-www-form-urlencoded, multipart/form-data, text/plain
	// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Access-Control-Expose-Headers
	AccessControlExposeHeaders []string

	// AccessControlAllowMethods (非简单请求的响应)所有允许请求跨域的方法，默认为 ["HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"]
	AccessControlAllowMethods []string

	// AccessControlAllowHeaders （非简单请求的响应）允许出默认之外，允许发送到服务端的 HEADER
	// https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/Access-Control-Allow-Headers
	AccessControlAllowHeaders []string

	// AccessControlMaxAge (非简单请求的响应)设置预检有效期，在有效期内，不需要再次发送预检请求
	// 单位：秒，不同的浏览器有上限
	AccessControlMaxAge string

	// accessControlAllOrigin 是否不限制任何来源，当 AccessControlAllowOrigin 含有 "*" 时为 true
	accessControlAllOrigin bool

	// accessControlAllowOrigin 将 AccessControlAllowOrigin 解析后存储于此
	accessControlAllowOrigin []*regexp.Regexp

	// accessControlExposeHeaders 将 AccessControlExposeHeaders 转为字符串形式存储，使用 "," 做为分隔符
	accessControlExposeHeaders string

	// accessControlAllowMethods 将 AccessControlAllowMethods 转为字符串形式存储，使用 "," 做为分隔符
	accessControlAllowMethods string

	// accessControlAllowHeaders 将 AccessControlAllowHeaders 转为字符串形式存储，使用 "," 做为分隔符
	accessControlAllowHeaders string
}

var defaultConfig = &Config{
	AccessControlAllowOrigin:      []string{"*"},
	AccessControlAllowCredentials: false,
	AccessControlExposeHeaders:    []string{},
	AccessControlAllowMethods: []string{
		http.MethodHead,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	},
	AccessControlAllowHeaders: []string{
		"Origin",
		"Accept",
		"Accept-Language",
		"Content-Language",
		"Content-Type",
	},
	AccessControlMaxAge: "600",
}

// New 跨域控制
//
// config 为 nil 时使用默认设置
func New(config *Config) zeroapi.Handler {
	c := defaultConfig
	c.init(config)

	return func(ctx zeroapi.Context) {
		// 预检请求
		if ctx.Method() == zeroapi.MethodOptions && len(ctx.Header("Access-Control-Request-Method")) != 0 {
			if c.checkPreflight(ctx) {
				// 预检通过
				c.checkPreflightSuccess(ctx)
			} else {
				// 预检未通过
				c.checkPreflightFailed(ctx)
			}

			return
		}

		// 浏览器的正常请求
		if c.checkRequest(ctx) {
			// 检查通过
			c.checkRequestSuccess(ctx)
		} else {
			// 检查未通过
			c.checkRequestFailed(ctx)
		}
	}
}

// checkPreflight 预检检查
func (c *Config) checkPreflight(ctx zeroapi.Context) bool {

	if !c.checkOrigin(ctx) {
		return false
	}

	method := ctx.Header("Access-Control-Request-Method")
	if !c.checkMethod(method) {
		return false
	}

	if !c.checkHeader(ctx) {
		return false
	}

	return true
}

func (c *Config) checkPreflightSuccess(ctx zeroapi.Context) {
	if c.accessControlAllOrigin && !c.AccessControlAllowCredentials {
		// Cookie 遵循同源政策，如果允许携带 Cooke，则不允许设置 Origin 为 *
		ctx.AddHeader("Access-Control-Allow-Origin", "*")
	} else {
		ctx.AddHeader("Access-Control-Allow-Origin", ctx.Header("Origin"))
	}

	ctx.AddHeader("Access-Control-Allow-Methods", c.accessControlAllowMethods)

	headers := ctx.Header("Access-Control-Request-Headers")
	if len(headers) > 0 {
		ctx.AddHeader("Access-Control-Allow-Headers", c.accessControlAllowHeaders)
	}

	if c.AccessControlAllowCredentials {
		ctx.AddHeader("Access-Control-Allow-Credentials", "true")
	}

	if c.AccessControlMaxAge != "" {
		ctx.AddHeader("Access-Control-Max-Age", c.AccessControlMaxAge)
	}

	ctx.Stopped()
	ctx.SetHTTPCode(http.StatusNoContent)
}

func (c *Config) checkPreflightFailed(ctx zeroapi.Context) {
	ctx.Stopped()
	ctx.SetHTTPCode(http.StatusForbidden)
	ctx.App().Logger().Errorf("preflight failed, origin: %s, method: %s", ctx.Header("Origin"), ctx.Method())
}

func (c *Config) checkRequest(ctx zeroapi.Context) bool {

	if !c.checkOrigin(ctx) {
		return false
	}

	method := ctx.Method()
	return c.checkMethod(method)
}

func (c *Config) checkRequestSuccess(ctx zeroapi.Context) {
	if c.accessControlAllOrigin && !c.AccessControlAllowCredentials {
		// Cookie 遵循同源政策，如果允许携带 Cooke，则不允许设置 Origin 为 *
		ctx.AddHeader("Access-Control-Allow-Origin", "*")
	} else {
		ctx.AddHeader("Access-Control-Allow-Origin", ctx.Header("Origin"))
	}

	if c.AccessControlAllowCredentials {
		ctx.AddHeader("Access-Control-Allow-Credentials", "true")
	}

	if len(c.accessControlExposeHeaders) > 0 {
		ctx.AddHeader("Access-Control-Expose-Headers", c.accessControlExposeHeaders)
	}
}

func (c *Config) checkRequestFailed(ctx zeroapi.Context) {
	ctx.SetHTTPCode(http.StatusForbidden)
	ctx.App().Logger().Errorf("request failed, origin: %s, method: %s", ctx.Header("Origin"), ctx.Method())
	ctx.Stopped()
}

func (c *Config) init(config *Config) {
	if config != nil {
		if len(config.AccessControlAllowOrigin) != 0 {
			c.AccessControlAllowOrigin = config.AccessControlAllowOrigin
		}
		c.AccessControlAllowCredentials = config.AccessControlAllowCredentials
		if len(config.AccessControlExposeHeaders) != 0 {
			c.AccessControlExposeHeaders = config.AccessControlExposeHeaders
		}
		if len(config.AccessControlAllowMethods) != 0 {
			c.AccessControlAllowMethods = config.AccessControlAllowMethods
		}
		if len(config.AccessControlAllowHeaders) != 0 {
			c.AccessControlAllowHeaders = config.AccessControlAllowHeaders
		}
		if config.AccessControlMaxAge != "" {
			c.AccessControlMaxAge = config.AccessControlMaxAge
		}
	}

	for _, origin := range c.AccessControlAllowOrigin {
		if origin == "*" {
			c.accessControlAllOrigin = true
			continue
		}

		pattern := regexp.QuoteMeta(origin)
		pattern = strings.Replace(pattern, "\\*", ".*", -1)
		pattern = strings.Replace(pattern, "\\?", ".", -1)
		p := "^" + pattern + "$"
		c.accessControlAllowOrigin = append(c.accessControlAllowOrigin, regexp.MustCompile(p))
	}

	c.accessControlExposeHeaders = strings.Join(c.AccessControlExposeHeaders, ",")
	c.accessControlAllowMethods = strings.Join(c.AccessControlAllowMethods, ",")
	c.accessControlAllowHeaders = strings.Join(c.AccessControlAllowHeaders, ",")
}

// checkOrigin 检查 ORIGIN
func (c *Config) checkOrigin(ctx zeroapi.Context) bool {
	if c.accessControlAllOrigin {
		return true
	}

	origin := ctx.Header("Origin")
	if origin == "" {
		return false
	}

	for _, r := range c.accessControlAllowOrigin {
		if r.MatchString(origin) {
			return true
		}
	}

	return false
}

// checkMethod 检查 Method
func (c *Config) checkMethod(method string) bool {
	if len(method) == 0 {
		return false
	}

	m := strings.ToUpper(method)

	for _, v := range c.AccessControlAllowMethods {
		if m == v {
			return true
		}
	}

	return false
}

// checkHeader 检查 HEADER
func (c *Config) checkHeader(ctx zeroapi.Context) bool {
	s := ctx.Header("Access-Control-Request-Headers")
	if len(s) == 0 {
		return true
	}
	headers := strings.Split(s, ",")
	if len(headers) == 0 {
		return true
	}

	keys := make([]string, 0, len(headers))
	for _, header := range headers {
		// http.CanonicalHeaderKey: 返回header key的规范话形式
		// 规范化形式是以"-"为分隔符，每一部分都是首字母大写，其他字母小写
		// 例如"accept-encoding" 的标准化形式是 "Accept-Encoding"

		key := http.CanonicalHeaderKey(header)
		keys = append(keys, key)
	}

	for _, key := range keys {
		found := false

		for _, v := range c.AccessControlAllowHeaders {

			if key == v {
				found = true
				break
			}
		}

		// 有一个不允许
		if !found {
			return false
		}
	}

	return true
}
