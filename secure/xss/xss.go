package xss

import (
	zeroapi "github.com/zerogo-hub/zero-api"
)

// New 添加信息头 X-XSS-Protection: 1; mode=block
// 启用XSS过滤。如果检测到攻击，浏览器将不会清除页面，而是阻止页面加载
func New() zeroapi.Handler {

	return func(ctx zeroapi.Context) {
		ctx.AddHeader("X-XSS-Protection", "1; mode=block")
	}
}
