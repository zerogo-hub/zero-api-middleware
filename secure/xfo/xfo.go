package xfo

import (
	zeroapi "github.com/zerogo-hub/zero-api"
)

// New 网站可以使用此功能，来确保自己网站的内容没有被嵌到别人的网站中去，也从而避免了点击劫持 (clickjacking) 的攻击
// enable: true 表示 X-Frame-Options: SAMEORIGIN，该页面可以在相同域名页面的 frame 中展示
// 		   false 表示 X-Frame-Options: DENY，表示该页面不允许在 frame 中展示，即便是在相同域名的页面中嵌套也不允许
// allowURI: 表示 X-Frame-Options: ALLOW-FROM uri，例如 X-Frame-Options: ALLOW-FROM https://keylala.cn
func New(enable bool, allowURI string) zeroapi.Handler {
	return func(ctx zeroapi.Context) {
		if allowURI != "" {
			ctx.AddHeader("X-Frame-Options", "ALLOW-FROM "+allowURI)
		} else if enable {
			ctx.AddHeader("X-Frame-Options", "SAMEORIGIN")
		} else {
			ctx.AddHeader("X-Frame-Options", "DENY")
		}
	}
}
